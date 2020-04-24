package html

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func FindNodes(resp *http.Response, selector string) *goquery.Selection {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse response body fail, err: %v", err)
	}

	return doc.Find(selector)
}
