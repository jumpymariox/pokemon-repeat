package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	file "reptile/internal/file"
	traverse "reptile/internal/html"
	http "reptile/internal/http"
	"reptile/internal/url"

	"golang.org/x/net/html"
)

const target = "https://wallpaperscraft.com"

func main() {
	queryFields := os.Args[1:]
	for _, field := range queryFields {
		targetURL := target + "/search/?query=" + field
		resBody, err := http.Fetch(targetURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request connect fail, url: %s\n: %v\n", field, err)
		} else {
			defer io.Copy(ioutil.Discard, resBody)

			parseHTML(targetURL, resBody)
		}
	}
}

func parseHTML(webURL string, body io.Reader) {
	doc, err := html.Parse(body)
	if err != nil {
		fmt.Errorf("parsing %s as HTML: %v\n ", webURL, err)
	}

	imgArray := []string{}
	imgArray = traverse.TraverseNodeAttr(doc, imgArray, "img", "src")

	for _, imgURL := range imgArray {
		resBody, err := http.Fetch(imgURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "download img fail, %s\n: %v\n", webURL, err)
		} else {
			defer io.Copy(ioutil.Discard, resBody)

			fileName, err := url.ParseFileName(imgURL)
			if err != nil {
				continue
			}
			file.Create(resBody, "output", fileName)
		}
	}
}
