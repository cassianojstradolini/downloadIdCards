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

// GenerateImageResponse contains the result of image generation
type GenerateImageResponse struct {
	ImageContent []byte
	FileName     string
}

// MergeImages creates a single image containing all ID card images from the response
func MergeImages(ctx context.Context, idCardsResp data.IdCardsResponseSchema) (*GenerateImageResponse, error) {
	var images []image.Image
	var pdfCards []data.IdCard

	// Process each card concurrently
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, card := range idCardsResp.Data {
		wg.Add(1)
		go func(card data.IdCard) {
			defer wg.Done()
			var img image.Image
			var err error

			if card.Attributes.Type == data.IdCardAttributesTypeHTML {
				mu.Lock()
				pdfCards = append(pdfCards, card)
				mu.Unlock()
				return
			} else {
				if isURL(card.Attributes.Source) {
					img, err = loadImageFromURL(card.Attributes.Source)
				} else {
					img, err = loadImageFromBase64(card.Attributes.Source)
				}
			}

			if err != nil {
				log.Printf("Failed to load image: %v", err)
				return
			}

			mu.Lock()
			images = append(images, img)
			mu.Unlock()
		}(card)
	}
	wg.Wait()

	// Process HTML cards separately
	if len(pdfCards) > 0 {
		htmlImages, err := ConvertHTMLCardsToImage(ctx, pdfCards)
		if err != nil {
			log.Printf("Warning: Failed to convert HTML cards: %v\n", err)
		} else {
			images = append(images, htmlImages...)
		}
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no valid images found to merge")
	}

	// Merge images vertically
	mergedImg, err := mergeImagesVertically(images)
	if err != nil {
		return nil, fmt.Errorf("failed to merge images: %w", err)
	}

	// Convert the merged image to bytes
	var buf bytes.Buffer
	if err := imgconv.Write(&buf, mergedImg, &imgconv.FormatOption{Format: imgconv.PNG}); err != nil {
		return nil, fmt.Errorf("failed to encode merged image: %w", err)
	}

	// Generate a filename
	fileName := fmt.Sprintf("id_cards_%s.png", time.Now().Format("20060102_150405"))

	return &GenerateImageResponse{
		ImageContent: buf.Bytes(),
		FileName:     fileName,
	}, nil
}

// we can resize HTML cards creating a merge specifically for HTML
func mergeImagesVertically(images []image.Image) (image.Image, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("no images to merge")
	}

	// Standard credit card size (in pixels at 300 DPI)
	cardWidth := 1012 // ~3.375 inches
	cardHeight := 638 // ~2.125 inches
	margin := 30      // Margin between cards
	sideMargin := 60  // Margin on sides

	// Create a new canvas with consistent width and calculated height
	totalHeight := (cardHeight * len(images)) + (margin * (len(images) - 1)) + (sideMargin * 2)
	totalWidth := cardWidth + (sideMargin * 2)
	mergedImg := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Fill background with white
	draw.Draw(mergedImg, mergedImg.Bounds(), image.White, image.Point{}, draw.Src)

	// Draw each image onto the canvas
	currentY := sideMargin
	for _, img := range images {
		// Get original dimensions
		bounds := img.Bounds()
		origWidth := bounds.Dx()
		origHeight := bounds.Dy()

		// Calculate scale factors for width and height
		widthRatio := float64(cardWidth) / float64(origWidth)
		heightRatio := float64(cardHeight) / float64(origHeight)

		// Use the smaller ratio to ensure image fits within card dimensions
		ratio := widthRatio
		if heightRatio < widthRatio {
			ratio = heightRatio
		}

		// Calculate new dimensions
		newWidth := int(float64(origWidth) * ratio)
		newHeight := int(float64(origHeight) * ratio)

		// Resize the image
		resizedImg := imgconv.Resize(img, &imgconv.ResizeOption{
			Width:  newWidth,
			Height: newHeight,
		})

		// Calculate x position to center the image
		xPos := sideMargin + ((cardWidth - newWidth) / 2)
		yPos := currentY + ((cardHeight - newHeight) / 2) // Also center vertically within its space

		// Draw the image
		draw.Draw(
			mergedImg,
			image.Rect(xPos, yPos, xPos+newWidth, yPos+newHeight),
			resizedImg,
			bounds.Min,
			draw.Over,
		)

		// Move to next card position
		currentY += cardHeight + margin
	}

	return mergedImg, nil
}

// ConvertHTMLCardsToImage converts HTML ID cards to a single image
func ConvertHTMLCardsToImage(ctx context.Context, htmlCards []data.IdCard) ([]image.Image, error) {
	// Create request with only HTML cards
	idCardsResp := data.IdCardsResponseSchema{
		Data: htmlCards,
	}

	// Generate PDF from HTML cards
	pdfResponse, err := to_pdf.GeneratePDFFromIDCards(ctx, idCardsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF from HTML cards: %w", err)
	}

	// Convert PDF to image
	img, err := convertPDFToImage(pdfResponse.PDFContent)
	if err != nil {
		return nil, fmt.Errorf("failed to convert PDF to image: %w", err)
	}

	// For now, we're returning a slice with a single image
	// If you need to separate each card into individual images, you would need
	// to modify the PDF generation or image processing to handle that
	return []image.Image{img}, nil
}

// convertPDFToImage converts PDF bytes to an image using unipdf
func convertPDFToImage(pdfBytes []byte) (image.Image, error) {
	log.Printf("Converting PDF to image using unipdf library")

	// Create a new PDF reader
	pdfReader, err := model.NewPdfReader(bytes.NewReader(pdfBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	// Get the first page
	page, err := pdfReader.GetPage(1)
	if err != nil {
		return nil, fmt.Errorf("failed to get PDF page: %w", err)
	}

	device := render.NewImageDevice()

	// Render the page to an image
	img, err := device.Render(page)
	if err != nil {
		return nil, fmt.Errorf("failed to render PDF page to image: %w", err)
	}

	return img, nil
}
