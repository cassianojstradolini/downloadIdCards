package data

import (
	"main/data"
)

type MockData struct {
	IdCards data.IdCardsResponseSchema
	Images  [][]byte // Raw image data for testing
}

// GenerateMockData creates a consistent set of test data
func GenerateMockData() MockData {
	return MockData{
		IdCards: data.IdCardsResponseSchema{
			Data: []data.IdCard{
				data.MockImageIdCardFront,
				data.MockImageIdCardBack,
				data.MockIdCardFront,
				data.MockIdCardBack,
				data.MockHTMLIdCardFront,
				data.MockHTMLIdCardBack,
				data.MockHTMLIdCardBoth,
			},
		},
	}
}
