package file

import (
	"fmt"
	"io"
	"os"
)

func Create(reader io.Reader, filePath string, fileName string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filePath, os.ModePerm)
	}
	file, createErr := os.Create(filePath + "/" + fileName)
	if createErr != nil {
		fmt.Println("create file fail", createErr)
	}
	_, copyErr := io.Copy(file, reader)
	if copyErr != nil {
		fmt.Println("copy reader fail", copyErr)
	}
}
