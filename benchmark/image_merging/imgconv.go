// image_merging/imgconv.go
package image_merging

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/sunshineplan/imgconv"
)

// MergeWithImgconv vertically merges a slice of images using imgconv for resizing
func MergeWithImgconv(images []image.Image) (image.Image, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("no images provided")
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
