package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	check "reptile/internal/errorcheck"
	"strings"
)

func Fetch(webURL string) (io.Reader, error) {
	if !checkURL(webURL) {
		err := fmt.Errorf("url is invalid", webURL)
		return nil, err
	}

	res, err := http.Get(webURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR fetch request %s:%v\n", webURL, err)
		return nil, err
	}
	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	check.Panic(err)

	return bytes.NewReader(bodyBytes), nil
}

func checkURL(str string) bool {
	return strings.HasPrefix(str, "https")
}
