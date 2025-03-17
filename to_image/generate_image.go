package to_image

import (
	"bytes"
	"context"
	"fmt"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
	"image"
	"image/draw"
	"log"
	"main/to_pdf"
	"sync"
	"time"

	"github.com/sunshineplan/imgconv"
	"main/data"
)

// Use sync.Pool for reusable buffers
var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GenerateImageResponse contains the result of image generation
type GenerateImageResponse struct {
	ImageContent []byte
	FileName     string
}

func MergeImages(ctx context.Context, idCardsResp data.IdCardsResponseSchema) (*GenerateImageResponse, error) {
	var images []image.Image
	var pdfCards []data.IdCard

	// Use a fixed number of workers for card processing.
	numWorkers := 10
	cardCh := make(chan data.IdCard)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Worker function for processing image cards.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for card := range cardCh {
				if card.Attributes.Type == data.IdCardAttributesTypeHTML {
					mu.Lock()
					pdfCards = append(pdfCards, card)
					mu.Unlock()
					continue
				}
				var img image.Image
				var err error
				if isURL(card.Attributes.Source) {
					img, err = loadImageFromURL(card.Attributes.Source)
				} else {
					img, err = loadImageFromBase64(card.Attributes.Source)
				}
				if err != nil {
					log.Printf("Failed to load image: %v", err)
					continue
				}
				mu.Lock()
				images = append(images, img)
				mu.Unlock()
			}
		}()
	}

	for _, card := range idCardsResp.Data {
		cardCh <- card
	}
	close(cardCh)
	wg.Wait()

	if len(pdfCards) > 0 {
		htmlImages, err := ConvertHTMLCardsToImage(ctx, pdfCards)
		if err != nil {
			log.Printf("Warning: Failed to convert HTML cards: %v", err)
		} else {
			images = append(images, htmlImages...)
		}
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no valid images found to merge")
	}

	mergedImg, err := mergeImagesVertically(images)
	if err != nil {
		return nil, fmt.Errorf("failed to merge images: %w", err)
	}

	// Get a buffer from the pool
	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)
	buf.Reset()

	// using JPEG due to lower memory footprint
	if err := imgconv.Write(buf, mergedImg, &imgconv.FormatOption{Format: imgconv.JPEG}); err != nil {
		return nil, fmt.Errorf("failed to encode merged image: %w", err)
	}

	fileName := fmt.Sprintf("id_cards_%s.png", time.Now().Format("20060102_150405"))
	return &GenerateImageResponse{
		ImageContent: buf.Bytes(),
		FileName:     fileName,
	}, nil
}

func mergeImagesVertically(images []image.Image) (image.Image, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("no images to merge")
	}

	cardWidth := 1012
	cardHeight := 638
	margin := 30
	sideMargin := 60

	totalHeight := (cardHeight * len(images)) + (margin * (len(images) - 1)) + (sideMargin * 2)
	totalWidth := cardWidth + (sideMargin * 2)
	mergedImg := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	draw.Draw(mergedImg, mergedImg.Bounds(), image.White, image.Point{}, draw.Src)

	currentY := sideMargin
	for _, img := range images {
		bounds := img.Bounds()
		origWidth := bounds.Dx()
		origHeight := bounds.Dy()

		widthRatio := float64(cardWidth) / float64(origWidth)
		heightRatio := float64(cardHeight) / float64(origHeight)
		ratio := widthRatio
		if heightRatio < widthRatio {
			ratio = heightRatio
		}

		newWidth := int(float64(origWidth) * ratio)
		newHeight := int(float64(origHeight) * ratio)

		resizedImg := imgconv.Resize(img, &imgconv.ResizeOption{
			Width:  newWidth,
			Height: newHeight,
		})

		xPos := sideMargin + ((cardWidth - newWidth) / 2)
		yPos := currentY + ((cardHeight - newHeight) / 2)

		draw.Draw(mergedImg, image.Rect(xPos, yPos, xPos+newWidth, yPos+newHeight), resizedImg, bounds.Min, draw.Over)
		currentY += cardHeight + margin
	}
	return mergedImg, nil
}

// Parallelize PDF generation in ConvertHTMLCardsToImage
func ConvertHTMLCardsToImage(ctx context.Context, htmlCards []data.IdCard) ([]image.Image, error) {
	numWorkers := 4
	cardCh := make(chan data.IdCard, len(htmlCards))
	imgCh := make(chan image.Image, len(htmlCards))
	errCh := make(chan error, len(htmlCards))

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for card := range cardCh {
				idCardsResp := data.IdCardsResponseSchema{Data: []data.IdCard{card}}
				pdfResponse, err := to_pdf.GeneratePDFFromIDCards(ctx, idCardsResp)
				if err != nil {
					errCh <- fmt.Errorf("failed to generate PDF from HTML card: %w", err)
					continue
				}
				img, err := convertPDFToImage(pdfResponse.PDFContent)
				if err != nil {
					errCh <- fmt.Errorf("failed to convert PDF to image: %w", err)
					continue
				}
				imgCh <- img
			}
		}()
	}

	for _, card := range htmlCards {
		cardCh <- card
	}
	close(cardCh)
	wg.Wait()
	close(imgCh)
	close(errCh)

	var images []image.Image
	for img := range imgCh {
		images = append(images, img)
	}

	for err := range errCh {
		log.Printf("Warning: %v\n", err)
	}
	return images, nil
}

func convertPDFToImage(pdfBytes []byte) (image.Image, error) {
	log.Printf("Converting PDF to image using unipdf library")
	pdfReader, err := model.NewPdfReader(bytes.NewReader(pdfBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}
	page, err := pdfReader.GetPage(1)
	if err != nil {
		return nil, fmt.Errorf("failed to get PDF page: %w", err)
	}
	device := render.NewImageDevice()
	img, err := device.Render(page)
	if err != nil {
		return nil, fmt.Errorf("failed to render PDF page to image: %w", err)
	}
	return img, nil
}
