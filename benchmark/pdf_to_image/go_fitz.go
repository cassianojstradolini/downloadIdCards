// pdf_to_image/go_fitz.go
package pdf_to_image

import (
	"fmt"
	"image"

	"github.com/gen2brain/go-fitz"
)

func ConvertWithGoFitz(pdfBytes []byte) (image.Image, error) {
	doc, err := fitz.NewFromMemory(pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}
	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		return nil, fmt.Errorf("failed to convert PDF page to image: %w", err)
	}

	return img, nil
}
