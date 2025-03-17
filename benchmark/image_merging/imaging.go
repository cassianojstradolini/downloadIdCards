// image_merging/imaging.go
package image_merging

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

func MergeWithImaging(images []image.Image) (image.Image, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("no images provided")
	}

	totalHeight := 0
	maxWidth := 0
	for _, img := range images {
		bounds := img.Bounds()
		totalHeight += bounds.Dy()
		if bounds.Dx() > maxWidth {
			maxWidth = bounds.Dx()
		}
	}

	dst := imaging.New(maxWidth, totalHeight, image.White)
	offsetY := 0
	for _, img := range images {
		dst = imaging.Paste(dst, img, image.Pt(0, offsetY))
		offsetY += img.Bounds().Dy()
	}

	return dst, nil
}
