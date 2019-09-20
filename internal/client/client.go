//Package client : Consume api V1 https://help.opendatasoft.com/apis/ods-search-v1/#search-api-v1
package client

import (
	"encoding/json"
	"log"
	"net/http"
)

const domain = "opendata.bordeaux-metropole.fr"
const endpointDatasets = "/api/datasets/1.0"
const endpointRecords = "/api/records/1.0"

// Api Quotas https://help.opendatasoft.com/apis/ods-search-v1/#quotas
type rateLimit struct {
	Limit     uint16
	Remaining uint16
	Reset     uint16
}

type Dataset struct {
	Datasetid string
}

func DatasetCatalog() {

	if err := fetchData("https://opendata.bordeaux-metropole.fr/api/datasets/1.0/search/"); err != nil {
		panic(err)
	}

}

func LookupDataset() {

}

func DownloadRecords() {

}

func LookupRecord() {

}

func fetchData(url string) error {

	// Get the data
	var v map[string]interface{}
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// @todo : throw error https://banzaicloud.com/blog/error-handling-go/
	log.Println(resp.Status)
	log.Println(resp.StatusCode)
	log.Println(resp.Header["X-Ratelimit-Limit"])
	log.Println(resp.Header["X-Ratelimit-Remaining"])
	log.Println(resp.Header["X-Ratelimit-Reset"])

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		log.Println(err)
		return err
	}

	//log.Println(v)

	return nil
}
