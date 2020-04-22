package html

import (
	"fmt"
	"io"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func FindNodes(reader io.Reader, selector string) *goquery.Selection {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse response body fail, err: %v", err)
	}

	return doc.Find(selector)
}
