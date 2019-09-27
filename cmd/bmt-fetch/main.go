//Package to fetch data from https://opendata.bordeaux-metropole.fr/api/datasets
package main

import (
	"BMTimetable/internal/client"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	filters := url.Values{}
	filters.Set("rows", "50")
	filters.Set("refine.keyword", "saeiv")

	datasets := client.DatasetCatalog(&filters)

	for _, dataset := range datasets {
		client.DownloadRecords(buildFilePath(dataset), dataset)
	}
}

func buildFilePath(dataset client.Dataset) string {

	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return filepath.Join(cwd, "tmp", dataset.Datasetid)
}
