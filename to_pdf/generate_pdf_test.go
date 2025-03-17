//go:build wkhtmltopdf

package to_pdf

import (
	"testing"

	"os"
	"strings"
)

func TestGeneratePDFFromIDCards(t *testing.T) {
	currentDir, _ := os.Getwd()

	type args struct {
		ctx         context.Context
		idCardsResp benefitsV2.IdCardsResponseSchema
	}
	tests := []struct {
		name    string
		args    args
		want    *GeneratePDFResponse
		wantErr bool
	}{
		{
			name: "Success with multiple ID cards",
			args: args{
				ctx: context.New(),
				idCardsResp: benefitsV2.IdCardsResponseSchema{
					Data: []benefitsV2.IdCard{
						mockBenefits.MockImageIdCardFront,
						mockBenefits.MockImageIdCardBack,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePDFFromIDCards(tt.args.ctx, tt.args.idCardsResp, currentDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePDFFromIDCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect success, verify file exists and has correct format
			if !tt.wantErr && got != nil {
				// Check if file exists
				if _, err := os.Stat(got.FilePath); os.IsNotExist(err) {
					t.Errorf("GeneratePDFFromIDCards() did not create file at %s", got.FilePath)
				}

				// Check filename format
				if !strings.HasPrefix(got.FileName, "id_cards_") || !strings.HasSuffix(got.FileName, ".pdf") {
					t.Errorf("GeneratePDFFromIDCards() generated incorrect filename format: %s", got.FileName)
				}

				// Clean up the file after test
				defer os.Remove(got.FilePath)
			}
		})
	}
}
