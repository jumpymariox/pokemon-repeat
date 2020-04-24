package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"reptile/internal/client"
	"reptile/internal/file"
	"reptile/internal/html"
	"reptile/internal/url"

	"github.com/PuerkitoBio/goquery"
)

const target = "https://wallpaperscraft.com"

func main() {
	// queryFields := os.Args[1:]
	queryFields := []string{"cat", "dog"}

	for _, field := range queryFields {
		targetURL := target + "/search/?query=" + field
		resp, err := client.Fetch(targetURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request connect fail, url: %s\n: %v\n", field, err)
		} else {
			download(targetURL, field, resp)
		}
	}
}

func download(webURL string, keyword string, resp *http.Response) {
	defer resp.Body.Close()
	html.FindNodes(resp, "a.wallpapers__link").Each(func(_ int, s *goquery.Selection) {
		href, isExist := s.Attr("href")
		if !isExist {
			return
		}

		segments := strings.Split(href, "/")

		downloadPageURL := target + "/download/" + segments[len(segments)-1] + "/1280x720"

		resp, err := client.Fetch(downloadPageURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "request connect fail, url: %s\n: %v\n", downloadPageURL, err)
			return
		}

		go downloadEachImg(webURL, keyword, resp)
	})
}

func downloadEachImg(webURL string, keyword string, resp *http.Response) {
	defer resp.Body.Close()
	html.FindNodes(resp, "a[download]").Each(func(_ int, s *goquery.Selection) {
		downloadURL, isExist := s.Attr("href")
		if !isExist {
			return
		}
		fmt.Println("downloading img", downloadURL)
		resp, err := client.Fetch(downloadURL)
		defer resp.Body.Close()
		if err != nil {
			fmt.Println("download img fail", downloadURL)
			fmt.Fprintf(os.Stderr, "download img fail, %s\n: %v\n", webURL, err)
		} else {
			fmt.Println("download img success", downloadURL)
			fileName, _ := url.ParseFileName(downloadURL)
			file.Create(resp.Body, "output/"+keyword, fileName)
		}
	})
}
