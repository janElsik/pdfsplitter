package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(filePath string, url string) error {
	fmt.Println("inside Download")
	//Get the data
	resp, err := http.Get(url)
	fmt.Println("Downloader:got url")

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("Downloader:defer body close")

	// Create the file

	out, err := os.Create(filePath)
	fmt.Println("Downloader: created blank file")

	if err != nil {
		return err
	}
	defer out.Close()
	fmt.Println("Downloader:defer out.close")

	// Write the body to file

	_, err = io.Copy(out, resp.Body)
	fmt.Println("Downloader:write body")

	return err
}
