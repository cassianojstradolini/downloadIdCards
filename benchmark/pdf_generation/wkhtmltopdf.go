// pdf_generation/wkhtmltopdf.go
package pdf_generation

import (
	"context"
	"main/data"
	"main/to_pdf"
)

// GenerateWithWkhtmltopdf generates a PDF using the wkhtmltopdf library
func GenerateWithWkhtmltopdf(idCards data.IdCardsResponseSchema) ([]byte, error) {
	resp, err := to_pdf.GeneratePDFFromIDCards(context.Background(), idCards)
	if err != nil {
		return nil, err
	}
	return resp.PDFContent, nil
}
