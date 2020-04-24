package client

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var httpclient http.Client

func main() {
	httpclient = http.Client{
		Timeout: 5 * time.Millisecond,
	}
}

func Fetch(webURL string) (*http.Response, error) {
	if !checkURL(webURL) {
		err := fmt.Errorf("url is invalid, url: %s", webURL)
		return nil, err
	}

	res, err := httpclient.Get(webURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR fetch request %s:%v\n", webURL, err)
		return nil, err
	}

	return res, err
}

func checkURL(str string) bool {
	return strings.HasPrefix(str, "https")
}
