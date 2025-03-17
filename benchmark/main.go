package main

import (
	"flag"
	"fmt"
	"log"
	"main/benchmark/data"
	"main/benchmark/image_merging"
	"main/benchmark/pdf_generation"
	"main/benchmark/pdf_to_image"
	"runtime"
	"sort"
	"time"
)

var (
	runPdfGen     = flag.Bool("pdf-gen", true, "Run PDF generation benchmarks")
	runPdfToImage = flag.Bool("pdf-to-img", true, "Run PDF to image conversion benchmarks")
	runImgMerge   = flag.Bool("img-merge", true, "Run image merging benchmarks")
	iterations    = flag.Int("n", 5, "Number of iterations for each benchmark")
	verbose       = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	fmt.Printf("=== PDF/Image Processing Benchmark ===\n")
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("CPU: %d cores\n\n", runtime.NumCPU())

	// Generate mock data
	mockData := data.GenerateMockData()

	// Run benchmarks
	if *runPdfGen {
		runPdfGenerationBenchmarks(mockData)
	}

	if *runPdfToImage {
		runPdfToImageBenchmarks()
	}

	if *runImgMerge {
		runImageMergingBenchmarks()
	}

	// Print summary of all benchmarks if all were run
	if *runPdfGen && *runPdfToImage && *runImgMerge {
		fmt.Println("\n=== BENCHMARK SUMMARY ===")
		fmt.Println("Based on the benchmark results, here are the recommended libraries for each operation:")
		fmt.Println()
	}
}

func runPdfGenerationBenchmarks(mockData data.MockData) {
	fmt.Println("=== PDF Generation Benchmarks ===")

	results := make(map[string]benchmarkResult)

	// wkhtmltopdf benchmark
	fmt.Println("- wkhtmltopdf:")
	results["wkhtmltopdf"] = benchmark(func() ([]byte, error) {
		return pdf_generation.GenerateWithWkhtmltopdf(mockData.IdCards)
	})
	printResults(results["wkhtmltopdf"])

	// chromedp benchmark (if available)
	fmt.Println("- chromedp:")
	results["chromedp"] = benchmark(func() ([]byte, error) {
		return pdf_generation.GenerateWithChromedp(mockData.IdCards)
	})
	printResults(results["chromedp"])

	// Compare results
	compareResults("PDF Generation", results)
	fmt.Println()
}

func runPdfToImageBenchmarks() {
	// First generate a sample PDF to convert
	mockData := data.GenerateMockData()
	pdfBytes, err := pdf_generation.GenerateWithWkhtmltopdf(mockData.IdCards)
	if err != nil {
		log.Fatalf("Failed to generate sample PDF: %v", err)
	}

	fmt.Println("=== PDF to Image Conversion Benchmarks ===")

	results := make(map[string]benchmarkResult)

	// unipdf benchmark
	fmt.Println("- unipdf:")
	results["unipdf"] = benchmark(func() ([]byte, error) {
		img, err := pdf_to_image.ConvertWithUnipdf(pdfBytes)
		if err != nil {
			return nil, err
		}
		return pdf_to_image.EncodeImage(img)
	})
	printResults(results["unipdf"])

	// go-fitz benchmark
	fmt.Println("- go-fitz (MuPDF):")
	results["go-fitz"] = benchmark(func() ([]byte, error) {
		img, err := pdf_to_image.ConvertWithGoFitz(pdfBytes)
		if err != nil {
			return nil, err
		}
		return pdf_to_image.EncodeImage(img)
	})
	printResults(results["go-fitz"])

	// Compare results
	compareResults("PDF to Image Conversion", results)
	fmt.Println()
}

func runImageMergingBenchmarks() {
	// First generate images to merge
	mockData := data.GenerateMockData()
	pdfBytes, err := pdf_generation.GenerateWithWkhtmltopdf(mockData.IdCards)
	if err != nil {
		log.Fatalf("Failed to generate sample PDF: %v", err)
	}

	// Convert PDF pages to images
	images, err := pdf_to_image.ExtractPagesToImages(pdfBytes)
	if err != nil {
		log.Fatalf("Failed to extract PDF pages: %v", err)
	}

	fmt.Println("=== Image Merging Benchmarks ===")

	results := make(map[string]benchmarkResult)

	// imaging benchmark
	fmt.Println("- imaging:")
	results["imaging"] = benchmark(func() ([]byte, error) {
		img, err := image_merging.MergeWithImaging(images)
		if err != nil {
			return nil, err
		}
		return pdf_to_image.EncodeImage(img)
	})
	printResults(results["imaging"])

	// imgconv benchmark
	fmt.Println("- imgconv:")
	results["imgconv"] = benchmark(func() ([]byte, error) {
		img, err := image_merging.MergeWithImgconv(images)
		if err != nil {
			return nil, err
		}
		return pdf_to_image.EncodeImage(img)
	})
	printResults(results["imgconv"])

	// standard library draw benchmark
	fmt.Println("- standard library (draw):")
	results["standard draw"] = benchmark(func() ([]byte, error) {
		img, err := image_merging.MergeWithDraw(images)
		if err != nil {
			return nil, err
		}
		return pdf_to_image.EncodeImage(img)
	})
	printResults(results["standard draw"])

	// Compare results
	compareResults("Image Merging", results)
	fmt.Println()
}

type benchmarkResult struct {
	iterations int
	duration   time.Duration
	bytesSize  int
	memStats   runtime.MemStats
	errors     []error
}

func benchmark(fn func() ([]byte, error)) benchmarkResult {
	var totalBytes int
	var memStats runtime.MemStats
	var errs []error

	runtime.GC()
	runtime.ReadMemStats(&memStats)
	beforeAlloc := memStats.TotalAlloc

	start := time.Now()
	for i := 0; i < *iterations; i++ {
		bytes, err := fn()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		totalBytes += len(bytes)
	}
	duration := time.Since(start)

	runtime.ReadMemStats(&memStats)
	memStats.TotalAlloc -= beforeAlloc

	return benchmarkResult{
		iterations: *iterations,
		duration:   duration,
		bytesSize:  totalBytes / *iterations,
		memStats:   memStats,
		errors:     errs,
	}
}

func printResults(r benchmarkResult) {
	fmt.Printf("  Time: %v (%v per operation)\n", r.duration, r.duration/time.Duration(r.iterations))
	fmt.Printf("  Memory: %.2f MB allocated\n", float64(r.memStats.TotalAlloc)/1024/1024)
	fmt.Printf("  Output size: %.2f KB\n", float64(r.bytesSize)/1024)

	if len(r.errors) > 0 {
		fmt.Printf("  Errors: %d of %d runs failed\n", len(r.errors), r.iterations)
		if *verbose {
			for i, err := range r.errors {
				fmt.Printf("    Error %d: %v\n", i+1, err)
			}
		}
	} else {
		fmt.Printf("  Errors: None\n")
	}
}

func compareResults(category string, results map[string]benchmarkResult) {
	if len(results) < 2 {
		return // No point comparing if there's only one result
	}

	type namedResult struct {
		name   string
		result benchmarkResult
	}

	var namedResults []namedResult
	for name, result := range results {
		namedResults = append(namedResults, namedResult{name: name, result: result})
	}

	fmt.Printf("\n=== %s Comparison ===\n", category)

	// 1. Performance comparison (speed)
	sort.Slice(namedResults, func(i, j int) bool {
		return namedResults[i].result.duration < namedResults[j].result.duration
	})
	fastest := namedResults[0]
	fmt.Printf("Speed ranking:\n")
	for i, r := range namedResults {
		opsPerSec := float64(r.result.iterations) / r.result.duration.Seconds()
		if i == 0 {
			fmt.Printf("  1. %s: %.2f ops/sec (fastest)\n", r.name, opsPerSec)
		} else {
			pctSlower := float64(r.result.duration-fastest.result.duration) / float64(fastest.result.duration) * 100
			fmt.Printf("  %d. %s: %.2f ops/sec (%.1f%% slower)\n", i+1, r.name, opsPerSec, pctSlower)
		}
	}

	// 2. Memory utilization comparison
	sort.Slice(namedResults, func(i, j int) bool {
		return namedResults[i].result.memStats.TotalAlloc < namedResults[j].result.memStats.TotalAlloc
	})
	leanest := namedResults[0]
	fmt.Printf("\nMemory usage ranking:\n")
	for i, r := range namedResults {
		memMB := float64(r.result.memStats.TotalAlloc) / 1024 / 1024
		if i == 0 {
			fmt.Printf("  1. %s: %.2f MB (most efficient)\n", r.name, memMB)
		} else {
			pctMore := float64(r.result.memStats.TotalAlloc-leanest.result.memStats.TotalAlloc) / float64(leanest.result.memStats.TotalAlloc) * 100
			fmt.Printf("  %d. %s: %.2f MB (%.1f%% more memory)\n", i+1, r.name, memMB, pctMore)
		}
	}

	// 3. Output size comparison
	sort.Slice(namedResults, func(i, j int) bool {
		return namedResults[i].result.bytesSize < namedResults[j].result.bytesSize
	})
	smallest := namedResults[0]
	fmt.Printf("\nOutput size ranking:\n")
	for i, r := range namedResults {
		sizeKB := float64(r.result.bytesSize) / 1024
		if i == 0 {
			fmt.Printf("  1. %s: %.2f KB (smallest)\n", r.name, sizeKB)
		} else {
			pctLarger := float64(r.result.bytesSize-smallest.result.bytesSize) / float64(smallest.result.bytesSize) * 100
			fmt.Printf("  %d. %s: %.2f KB (%.1f%% larger)\n", i+1, r.name, sizeKB, pctLarger)
		}
	}

	// 4. Overall recommendation section
	fmt.Printf("\nRecommendation for %s:\n", category)

	if category == "PDF to Image Conversion" {
		fmt.Printf("Best libraries by metric:\n")
		fmt.Printf("  Speed: %s\n", fastest.name)
		fmt.Printf("  Memory: %s\n", leanest.name)
		fmt.Printf("  Output size: %s\n", smallest.name)

		fmt.Printf("\nDetailed comparisons:\n")
		// For each library being compared
		for _, lib := range namedResults {
			// If it's the fastest, get comparison metrics against other libraries
			if lib.name == fastest.name {
				fmt.Printf("%s (fastest):\n", lib.name)

				// Compare against other libraries
				for _, other := range namedResults {
					if other.name == lib.name {
						continue // Skip self comparison
					}

					// Calculate how much faster this lib is compared to others
					speedAdvantage := float64(other.result.duration-lib.result.duration) / float64(lib.result.duration) * 100

					// Calculate memory usage differences
					var memText string
					if lib.name == leanest.name {
						memAdvantage := float64(other.result.memStats.TotalAlloc-lib.result.memStats.TotalAlloc) / float64(lib.result.memStats.TotalAlloc) * 100
						memText = fmt.Sprintf("uses %.1f%% less memory", memAdvantage)
					} else {
						memDisadvantage := float64(lib.result.memStats.TotalAlloc-leanest.result.memStats.TotalAlloc) / float64(leanest.result.memStats.TotalAlloc) * 100
						memText = fmt.Sprintf("uses %.1f%% more memory than %s", memDisadvantage, leanest.name)
					}

					// Calculate output size differences
					var sizeText string
					if lib.name == smallest.name {
						sizeAdvantage := float64(other.result.bytesSize-lib.result.bytesSize) / float64(lib.result.bytesSize) * 100
						sizeText = fmt.Sprintf("produces %.1f%% smaller output", sizeAdvantage)
					} else {
						sizeDisadvantage := float64(lib.result.bytesSize-smallest.result.bytesSize) / float64(smallest.result.bytesSize) * 100
						sizeText = fmt.Sprintf("produces %.1f%% larger output than %s", sizeDisadvantage, smallest.name)
					}

					fmt.Printf(" -  vs %s: %.1f%% faster, %s, %s\n", other.name, speedAdvantage, memText, sizeText)
				}
			}

			// If it's not the fastest but has the smallest output, highlight that
			if lib.name != fastest.name && lib.name == smallest.name {
				fmt.Printf("%s (smallest output):\n", lib.name)
				speedDisadvantage := float64(lib.result.duration-fastest.result.duration) / float64(fastest.result.duration) * 100
				memoryInfo := ""
				if lib.name == leanest.name {
					memoryInfo = "most memory efficient, "
				} else {
					memDiff := float64(lib.result.memStats.TotalAlloc-leanest.result.memStats.TotalAlloc) / float64(leanest.result.memStats.TotalAlloc) * 100
					memoryInfo = fmt.Sprintf("uses %.1f%% more memory than %s, ", memDiff, leanest.name)
				}
				fmt.Printf(" -  %.1f%% slower than %s, %sbut produces the smallest output\n",
					speedDisadvantage, fastest.name, memoryInfo)
			}
		}
	} else if category == "PDF Generation" || category == "Image Merging" {
		// For other categories, provide a more straightforward comparison
		if fastest.name != smallest.name {
			// If fastest is not also smallest output
			fmt.Printf("  %s is the fastest option:\n", fastest.name)

			// Find the second fastest for comparison
			var secondFastest namedResult
			for _, nr := range namedResults {
				if nr.name != fastest.name {
					secondFastest = nr
					break
				}
			}

			speedAdvantage := float64(secondFastest.result.duration-fastest.result.duration) /
				float64(secondFastest.result.duration) * 100

			fmt.Printf(" -    %.1f%% faster than %s\n", speedAdvantage, secondFastest.name)

			// Memory comparison
			if fastest.name != leanest.name {
				memDisadvantage := float64(fastest.result.memStats.TotalAlloc-leanest.result.memStats.TotalAlloc) /
					float64(leanest.result.memStats.TotalAlloc) * 100
				fmt.Printf(" -    Uses %.1f%% more memory than %s\n", memDisadvantage, leanest.name)
			} else {
				fmt.Printf(" -    Also most memory efficient\n")
			}

			// Output size comparison if not smallest
			sizeDisadvantage := float64(fastest.result.bytesSize-smallest.result.bytesSize) /
				float64(smallest.result.bytesSize) * 100
			if sizeDisadvantage > 0 {
				fmt.Printf(" -    Produces %.1f%% larger output than %s\n", sizeDisadvantage, smallest.name)
			}

			// Smallest output recommendation if different from fastest
			var memoryText string
			if smallest.name == leanest.name {
				memoryText = " and uses least memory"
			} else {
				memoryText = ""
			}
			fmt.Printf("\n  %s produces the smallest output size (%.2f KB)%s\n",
				smallest.name, float64(smallest.result.bytesSize)/1024, memoryText)
		} else {
			// If fastest is also smallest
			fmt.Printf("  %s is the fastest option and produces the smallest output.\n", fastest.name)
			if fastest.name != leanest.name {
				memDisadvantage := float64(fastest.result.memStats.TotalAlloc-leanest.result.memStats.TotalAlloc) /
					float64(leanest.result.memStats.TotalAlloc) * 100
				fmt.Printf("  However, it uses %.1f%% more memory than %s.\n", memDisadvantage, leanest.name)
			} else {
				fmt.Printf("  It is also the most memory efficient option.\n")
			}
		}
	} else {
		fmt.Printf("  %s is recommended based on performance metrics.\n", fastest.name)
	}
}
