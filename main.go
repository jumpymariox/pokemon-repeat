package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	check "reptile/internal/errorcheck"
	file "reptile/internal/file"
	traverse "reptile/internal/html"
	"reptile/internal/http"

	"golang.org/x/net/html"
)

func main() {
	webURLArray := [1]string{"https://wallpaperscraft.com/tag/cat/1920x1080/"}
	for _, webURL := range webURLArray {
		resBody, err := http.Fetch(webURL)
		if err != nil {
			fmt.Println("request connect fail, url:" + webURL)
		} else {
			defer io.Copy(ioutil.Discard, resBody)

			parseHTML(webURL, resBody)
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
			fmt.Println("download img fail, url:" + webURL)
		} else {
			defer io.Copy(ioutil.Discard, resBody)

			fileName := parseFileName(imgURL)
			file.Create(resBody, "output", fileName)
		}
	}
}

func parseFileName(fileURL string) string {
	parsedURL, err := url.Parse(fileURL)
	check.Panic(err)

	path := parsedURL.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]
	return fileName
}
