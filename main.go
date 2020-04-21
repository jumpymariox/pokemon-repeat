package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	check "reptile/internal/errorcheck"
	file "reptile/internal/file"
	selector "reptile/internal/html"

	"golang.org/x/net/html"
)

func main() {
	webURLArray := [1]string{"https://wallpaperscraft.com/tag/cat/1920x1080/"}
	for _, webURL := range webURLArray {
		resBody := fetch(webURL)
		parseHTML(webURL, resBody)
	}
}

func fetch(webURL string) io.Reader {
	if !checkURL(webURL) {
		fmt.Println("url is invalid", webURL)
		os.Exit(1)
	}

	res, err := http.Get(webURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR fetch request %s:%v\n", webURL, err)
		os.Exit(1)
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	check.Panic(err)

	return bytes.NewReader(bodyBytes)
}

func checkURL(str string) bool {
	return strings.HasPrefix(str, "https")
}

func parseHTML(webURL string, body io.Reader) {
	doc, err := html.Parse(body)
	if err != nil {
		fmt.Errorf("parsing %s as HTML: %v\n ", webURL, err)
	}

	imgArray := []string{}
	imgArray = selector.TraverseNodeAttr(doc, imgArray, "img", "src")

	for _, imgURL := range imgArray {
		if !checkURL(imgURL) {
			fmt.Println("url is invalid", imgURL)
			continue
		}

		res, err := http.Get(imgURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR fetch request %s:%v\n", imgURL, err)
			os.Exit(1)
		}
		defer res.Body.Close()
		defer io.Copy(ioutil.Discard, res.Body)

		fileName := parseFileName(imgURL)
		file.Create(res.Body, "output", fileName)
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
