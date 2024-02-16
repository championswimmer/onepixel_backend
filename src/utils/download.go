package utils

import (
	"fmt"
	"io"
	"net/http"
	"onepixel_backend/src/utils/applogger"
	"os"
	"time"
)

func IsFileFresh(maxdays int, filepath string) bool {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		applogger.Error("CheckFileAge: ", err)
		return false
	}
	fileAge := time.Since(fileInfo.ModTime())
	return fileAge.Hours() < float64(maxdays*24)
}

func DownloadFile(url string, filepath string) error {
	// Send a GET request to the file URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create a new file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	return err
}
