package main

import (
	"io"
	"net/http"
	"os"
)

func main() {

	fileUrl := "https://opendata.bordeaux-metropole.fr/api/datasets/1.0/search/?rows=100&start=0&refine.keyword=saeiv"

	if err := DownloadFile("datasets.json", fileUrl); err != nil {
		panic(err)
	}
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
