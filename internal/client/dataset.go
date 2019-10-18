//Package client : Download Dataset
package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

type Catalog struct {
	Nhits      uint
	Parameters Parameters
	Datasets   []Dataset
}

type Dataset struct {
	Datasetid string
	Metas     struct {
		Publisher         string
		Domain            string
		RecordsCount      uint `json:"records_count"`
		Title             string
		MetadataProcessed string `json:"metadata_processed"`
		DataProcessed     string `json:"data_processed"`
	}
	HasRecords bool
	fields     []interface{}
}

// https://blog.golang.org/pipelines

// Download pipeline
// 1rst stage : Get List of All datasets to dowload
// 2nd stage : Download File

// FetchDatasets Téléchargement parralèle des datasets
func FetchDatasets(parameters *url.Values) {

	catalog := &Catalog{}

	// Get Catalog
	FetchDatasetsCatalog(parameters, catalog)

	datasets := browseDatasets(catalog)

	files := downloadDatasets(datasets)

	log.Printf("Fin FetchDatasets %v", files)
}

//browseCatalog create a buffered channel and send into it
//all the catalogs to download
func browseDatasets(catalog *Catalog) <-chan Dataset {
	nbCatalogs := len(catalog.Datasets)
	c := make(chan Dataset, nbCatalogs)

	for i := 0; i < nbCatalogs; i++ {
		log.Printf("dataset %d to DL : %v", i, catalog.Datasets[i].Datasetid)
		c <- catalog.Datasets[i]
	}

	close(c)

	return c
}

func downloadDatasets(datasets <-chan Dataset) []string {
	// close done channel even if all downloading are not finished
	done := make(chan struct{})
	defer close(done)

	// Start a fixed number of goroutines to download files
	c := make(chan result) // HLc
	var wg sync.WaitGroup
	const numClients = 4
	wg.Add(numClients)
	for i := 0; i < numClients; i++ {
		go func() {
			downloadDataset(done, datasets, c) // HLc
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c) // HLc
	}()

	//sink
	var files []string
	for r := range c {
		if r.err != nil {
			log.Printf("File %v has not been downloaded : %v", r.dest, r.err)
		} else {
			files = append(files, r.dest)
			log.Printf("Sink Download completed %v : %v", len(files), r.dest)
		}
	}

	return files
}

type result struct {
	dest    string
	dataset Dataset
	err     error
}

func downloadDataset(done <-chan struct{}, datasets <-chan Dataset, files chan<- result) {

	for dataset := range datasets { // HLpaths

		jsonFile, err := path(dataset, "json")
		if err != nil {
			log.Println(err)
			continue
		}

		b, err := json.Marshal(dataset)
		if err != nil {
			log.Println(err)
			continue
		}

		if err := ioutil.WriteFile(jsonFile, b, os.ModePerm); err != nil {
			log.Println(err)
			continue
		}

		csvFile, err := path(dataset, "csv")
		if err != nil {
			log.Println(err)
			continue
		}

		select {
		case files <- result{csvFile, dataset, DownloadDataset(dataset, csvFile)}:
		case <-done:
			log.Println("WtF done closed !")
			return
		}
	}
}

func path(dataset Dataset, ext string) (string, error) {
	dir, err := directory(dataset)

	if err != nil {
		return "", err
	}

	return filepath.Join(dir, dataset.Datasetid+"."+ext), nil
}

func directory(dataset Dataset) (string, error) {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	dir := filepath.Join(cwd, "tmp", dataset.Datasetid)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	return dir, nil
}
