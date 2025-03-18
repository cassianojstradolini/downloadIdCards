# PDF/Image Processing Benchmark

This project benchmarks different libraries for PDF generation, PDF-to-image conversion, and image merging operations. It provides detailed comparisons of performance metrics including speed, memory usage, and output size.

## Features

- **PDF Generation**: Compare `wkhtmltopdf` and `chromedp`
- **PDF to Image Conversion**: Compare `go-fitz` (MuPDF) and `unipdf`
- **Image Merging**: Compare Go standard library (`draw`), `imgconv`, and `imaging`
- **Comprehensive Benchmarking**: Measure operations per second, memory allocation, and output size

## Installation

```bash
git clone https://github.com/cassianojstradolini/downloadIdCards.git
cd downloadIdCards
go mod download
```

## Usage

### Running the HTTP Server

The project includes a server that exposes endpoints for PDF and image operations:

```bash
go run main.go
```

Server runs on port 8081 with the following endpoints:
- `/pdf/idcards`: Generate PDF from ID cards
- `/image/idcards`: Generate merged image from ID cards
- `/template-extension/idcards`: Template extension functionality

### Running Benchmarks

```bash
go run benchmark/main.go [flags]
```

Available flags:
- `-pdf-gen=true|false`: Run PDF generation benchmarks (default: true)
- `-pdf-to-img=true|false`: Run PDF to image conversion benchmarks (default: true)
- `-img-merge=true|false`: Run image merging benchmarks (default: true)
- `-n=<number>`: Set number of iterations for each benchmark (default: 5)
- `-v`: Enable verbose output (default: false)

## Project Structure

- `main.go` - HTTP server for PDF/image generation
- `to_pdf/` - PDF generation implementation
- `to_image/` - Image generation & merging implementation
- `data/` - Mock data for testing
- `benchmark/` - Benchmark implementations:
  - `main.go` - Benchmark runner
  - `pdf_generation/` - PDF generation benchmarks
  - `pdf_to_image/` - PDF to image conversion benchmarks
  - `image_merging/` - Image merging benchmarks

## Dependencies

- **PDF Generation**
  - `github.com/SebastiaanKlippert/go-wkhtmltopdf` - HTML to PDF conversion
  - `github.com/chromedp/chromedp` - Browser automation for PDF generation

- **PDF to Image**
  - `github.com/gen2brain/go-fitz` - MuPDF-based PDF rendering
  - `github.com/unidoc/unipdf/v3` - PDF processing library

- **Image Processing**
  - `golang.org/x/image/draw` - Go standard library drawing
  - `github.com/disintegration/imaging` - Image processing utilities
  - `github.com/sunshineplan/imgconv` - Image conversion

## Benchmark Output

The benchmark provides detailed output including:
- Speed rankings (operations per second)
- Memory usage rankings (MB allocated)
- Output size rankings (KB)
- Detailed comparisons between libraries
- Recommendations based on different metrics

## Example

```
=== PDF Generation Comparison ===
Speed ranking:
  1. chromedp: 0.56 ops/sec (fastest)
  2. wkhtmltopdf: 0.42 ops/sec (31.5% slower)

Memory usage ranking:
  1. wkhtmltopdf: 10.76 MB (most efficient)
  2. chromedp: 186.56 MB (1634.2% more memory)

Output size ranking:
  1. wkhtmltopdf: 433.00 KB (smallest)
  2. chromedp: 763.47 KB (76.3% larger)

Recommendation for PDF Generation:
Best libraries by metric:
  Speed: chromedp
  Memory: wkhtmltopdf
  Output size: wkhtmltopdf

Detailed comparisons:
chromedp (fastest):
 - vs wkhtmltopdf: 23.9% faster, uses 1634.2% more memory than wkhtmltopdf, produces 76.3% larger output than wkhtmltopdf

wkhtmltopdf (smallest output):
 - 31.5% slower than chromedp, most memory efficient, but produces the smallest output

Summary: Choose chromedp for speed, wkhtmltopdf for memory efficiency and smallest output size.
```

## License

This project is licensed under the MIT License.
