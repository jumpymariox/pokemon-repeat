package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"reptile/internal/file"
	"reptile/internal/html"
	http "reptile/internal/http"
	"reptile/internal/url"

	"github.com/PuerkitoBio/goquery"
)

const target = "https://wallpaperscraft.com"

func main() {
	//queryFields := os.Args[1:]
	queryFields := [1]string{"test"}

	for _, field := range queryFields {
		targetURL := target + "/search/?query=" + field
		resBody, err := http.Fetch(targetURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request connect fail, url: %s\n: %v\n", field, err)
		} else {
			defer io.Copy(ioutil.Discard, resBody)

			downloadTargetImgs(targetURL, field, resBody)
		}
	}
}

func downloadTargetImgs(webURL string, keyword string, body io.Reader) {
	html.FindNodes(body, "a.wallpapers__link").Each(func(_ int, s *goquery.Selection) {
		href, isExist := s.Attr("href")
		if !isExist {
			return
		}

		segments := strings.Split(href, "/")

		downloadPageURL := target + "/download/" + segments[len(segments)-1] + "/1080x1920"
		resBody, err := http.Fetch(downloadPageURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request connect fail, url: %s\n: %v\n", downloadPageURL, err)
			return
		}
		html.FindNodes(resBody, "a[download]").Each(func(_ int, s *goquery.Selection) {
			downloadURL, isExist := s.Attr("href")
			if !isExist {
				return
			}
			fmt.Println("downing url", downloadURL)
			resBody, err := http.Fetch(downloadURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "download img fail, %s\n: %v\n", webURL, err)
			} else {
				defer io.Copy(ioutil.Discard, resBody)

				fileName, _ := url.ParseFileName(downloadURL)
				file.Create(resBody, "output/"+keyword, fileName)
			}
		})

	})
}
