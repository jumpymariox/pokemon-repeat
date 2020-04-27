package replite

import (
	"fmt"
	"os"
	"reptile/internal/tool/client"
	"reptile/internal/tool/file"
	"reptile/internal/tool/html"
	"reptile/internal/tool/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const target = "https://wallpaperscraft.com"

func Replite() {

	// queryFields := os.Args[1:]

	queryFields := []string{"cat"}
	imagePageURLChan := make(chan string, 30)
	downloadURLChan := make(chan string, 30)
	done := make(chan bool)

	go collectPageURL(imagePageURLChan, queryFields)
	go collectDownloadURL(imagePageURLChan, downloadURLChan)
	go downloadURL(downloadURLChan, done)

	<-done
}

func collectPageURL(urlChan chan<- string, queryFields []string) {
	defer func() {
		close(urlChan)
		fmt.Println("collect page url complete")
	}()

	for _, field := range queryFields {
		targetURL := target + "/search/?query=" + field

		resp, err := client.Fetch(targetURL)
		defer resp.Body.Close()
		if err != nil {
			continue
		}
		html.FindNodes(resp, "a.wallpapers__link").Each(func(_ int, s *goquery.Selection) {
			href, isExist := s.Attr("href")
			if !isExist {
				return
			}

			segments := strings.Split(href, "/")
			urlChan <- target + "/download/" + segments[len(segments)-1] + "/1280x720"
		})
	}
}

func collectDownloadURL(urlChan <-chan string, collectChan chan<- string) {
	defer func() {
		close(collectChan)
		fmt.Println("collect download url complete")
	}()
	for url := range urlChan {
		resp, err := client.Fetch(url)
		defer resp.Body.Close()
		if err != nil {
			continue
		}
		html.FindNodes(resp, "a[download]").Each(func(_ int, s *goquery.Selection) {
			downloadURL, isExist := s.Attr("href")
			if !isExist {
				return
			}
			collectChan <- downloadURL
		})
	}
}

func downloadURL(downloadURLChan <-chan string, done chan<- bool) {
	defer func() {
		fmt.Println("download all complete")
		done <- true
	}()
	for downloadURL := range downloadURLChan {
		resp, err := client.Fetch(downloadURL)
		defer func() {
			fmt.Sprintf("resp: %v/n", resp)
			resp.Body.Close()
		}()
		if err != nil {
			fmt.Fprintf(os.Stderr, "download img fail, %s\n: %v\n", downloadURL, err)
		} else {
			fmt.Println("download img success", downloadURL)
			fileName, _ := url.ParseFileName(downloadURL)
			file.Create(resp.Body, "output/", fileName)
		}
	}
}
