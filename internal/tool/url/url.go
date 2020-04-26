package url

import (
	"net/url"
	"strings"
)

func ParseFileName(fileURL string) (string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	path := parsedURL.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]
	return fileName, nil
}
