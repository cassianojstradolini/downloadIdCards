package to_pdf

import (
	"bytes"
	"context"
	"fmt"
	"main/data"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// GeneratePDFResponse contains the result of PDF generation
type GeneratePDFResponse struct {
	PDFContent []byte
	FileName   string
}

func GeneratePDFFromIDCards(ctx context.Context, idCardsResp data.IdCardsResponseSchema) (*GeneratePDFResponse, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize PDF generator: %w", err)
	}

	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeLetter)
	pdfg.MarginTop.Set(40)
	pdfg.MarginBottom.Set(40)
	pdfg.MarginLeft.Set(40)
	pdfg.MarginRight.Set(40)

	var sb strings.Builder
	for _, card := range idCardsResp.Data {
		if card.Attributes.Type == data.IdCardAttributesTypeHTML {
			sb.WriteString(`<div class="card">`)
			sb.WriteString(card.Attributes.Source)
			sb.WriteString(`</div>`)
		} else {
			var imgSrc string
			if isURL(card.Attributes.Source) {
				imgSrc = card.Attributes.Source
			} else {
				imgSrc = fmt.Sprintf("data:image/png;base64,%s", card.Attributes.Source)
			}
			sb.WriteString(fmt.Sprintf(
				`<div class="card"><img src="%s" alt="%s Card"></div>`,
				imgSrc, card.Attributes.Face))
		}
	}

	html := fmt.Sprintf(`
	  <!DOCTYPE html>
	  <html>
	  <head>
	   <style>
		body {
		 margin: 0;
		 padding: 0;
		 font-family: Arial, sans-serif;
		}
		.card {
		 width: 100%%;
		 margin-bottom: 20px;
		}
		img {
		 width: 100%%;
		 height: auto;
		}
		.card-info {
		 margin-top: 5px;
		 font-size: 12px;
		}
	   </style>
	  </head>
	  <body>
	   %s
	  </body>
	  </html>`, sb.String())

	page := wkhtmltopdf.NewPageReader(bytes.NewReader([]byte(html)))
	pdfg.AddPage(page)

	if err = pdfg.Create(); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	fileName := fmt.Sprintf("id_cards_%s.pdf", time.Now().Format("20060102_150405"))
	return &GeneratePDFResponse{
		PDFContent: pdfg.Bytes(),
		FileName:   fileName,
	}, nil
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
