// pdf_to_image/unipdf.go
package pdf_to_image

import (
	"bytes"
	"fmt"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
	"image"
	"image/jpeg"
)

// ConvertWithUnipdf converts a PDF to an image using the unipdf library
func ConvertWithUnipdf(pdfBytes []byte) (image.Image, error) {
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

// ExtractPagesToImages extracts all pages from a PDF and converts them to images
func ExtractPagesToImages(pdfBytes []byte) ([]image.Image, error) {
	pdfReader, err := model.NewPdfReader(bytes.NewReader(pdfBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get number of pages: %w", err)
	}

	var images []image.Image
	device := render.NewImageDevice()

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return nil, fmt.Errorf("failed to get page %d: %w", i, err)
		}

		img, err := device.Render(page)
		if err != nil {
			return nil, fmt.Errorf("failed to render page %d: %w", i, err)
		}

		images = append(images, img)
	}

	return images, nil
}

// EncodeImage encodes an image to JPEG format
func EncodeImage(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
