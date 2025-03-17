// image_merging/draw.go
package image_merging

import (
	"fmt"
	"golang.org/x/image/draw" // Extended draw package
	"image"
)

// MergeWithDraw vertically merges a slice of images using the extended draw package
func MergeWithDraw(images []image.Image) (image.Image, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("no images provided")
	}

	cardWidth := 1012
	cardHeight := 638
	margin := 30
	sideMargin := 60

	totalHeight := (cardHeight * len(images)) + (margin * (len(images) - 1)) + (sideMargin * 2)
	totalWidth := cardWidth + (sideMargin * 2)

	// Create a new RGBA image to hold the merged result
	mergedImg := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	// Fill with white background (use standard draw for this simple operation)
	draw.Draw(mergedImg, mergedImg.Bounds(), image.White, image.Point{}, draw.Src)

	currentY := sideMargin
	for _, img := range images {
		bounds := img.Bounds()
		origWidth := bounds.Dx()
		origHeight := bounds.Dy()

		// Calculate scaling ratios
		widthRatio := float64(cardWidth) / float64(origWidth)
		heightRatio := float64(cardHeight) / float64(origHeight)
		ratio := widthRatio
		if heightRatio < widthRatio {
			ratio = heightRatio
		}

		newWidth := int(float64(origWidth) * ratio)
		newHeight := int(float64(origHeight) * ratio)

		// Calculate position
		xPos := sideMargin + ((cardWidth - newWidth) / 2)
		yPos := currentY + ((cardHeight - newHeight) / 2)

		draw.NearestNeighbor.Scale(
			mergedImg,
			image.Rect(xPos, yPos, xPos+newWidth, yPos+newHeight),
			img,
			img.Bounds(),
			draw.Over,
			nil,
		)

		currentY += cardHeight + margin
	}

	return mergedImg, nil
}
