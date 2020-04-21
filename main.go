package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	file "reptile/internal/file"
	selector "reptile/internal/html"

	"golang.org/x/net/html"
)

func main() {
	urlArray := [1]string{"https://wallpaperscraft.com/tag/cat/1920x1080/"}
	for _, url := range urlArray {
		fetch(url)
	}
}

func fetch(url string) {
	if !checkURL(url) {
		fmt.Println("url is invalid", url)
		os.Exit(1)
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR fetch request %s:%v\n", url, err)
		os.Exit(1)
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	file.Create(res.Body, "output", "test.html")
	parseHTML(url, res.Body)
}

func checkURL(str string) bool {
	return strings.HasPrefix(str, "https")
}

func parseHTML(url string, body io.Reader) {
	doc, err := html.Parse(body)
	if err != nil {
		fmt.Errorf("parsing %s as HTML: %v\n ", url, err)
	}

	imgArray := []string{}
	imgArray = selector.TraverseNodeAttr(doc, imgArray, "img", "src")

}
