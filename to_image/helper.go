package to_image

import (
	"bytes"
	"encoding/base64"
	"github.com/sunshineplan/imgconv"
	"image"
	"net/http"
	"strings"
)

// Helper functions
func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func loadImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return imgconv.Decode(resp.Body)
}

func loadImageFromBase64(base64Str string) (image.Image, error) {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	return imgconv.Decode(bytes.NewReader(data))
}
