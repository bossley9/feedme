package api

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func FetchHTML(url string) (*goquery.Document, error) {
	res, err := FetchGet(url)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(bytes.NewReader(res))
}
