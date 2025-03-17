// pdf_generation/chromedp.go
package pdf_generation

import (
	"context"
	"fmt"
	"main/data"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// GenerateWithChromedp generates a PDF using the chromedp library
func GenerateWithChromedp(idCards data.IdCardsResponseSchema) ([]byte, error) {
	// Create a new Chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Generate HTML content
	var sb strings.Builder
	for _, card := range idCards.Data {
		if card.Attributes.Type == data.IdCardAttributesTypeHTML {
			sb.WriteString(`<div class="card">`)
			sb.WriteString(card.Attributes.Source)
			sb.WriteString(`</div>`)
		} else {
			var imgSrc string
			if strings.HasPrefix(card.Attributes.Source, "http") {
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

	// Navigate to data URL with our HTML content
	var pdfContent []byte
	dataURL := "data:text/html," + html
	err := chromedp.Run(ctx,
		chromedp.Navigate(dataURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfContent, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfContent, nil
}
