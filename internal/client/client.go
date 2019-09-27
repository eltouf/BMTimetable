//Package client : Consume api V1 https://help.opendatasoft.com/apis/ods-search-v1/#search-api-v1
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const host = "opendata.bordeaux-metropole.fr"

// Api Quotas https://help.opendatasoft.com/apis/ods-search-v1/#quotas
type RateLimitError struct {
	Limit     uint16
	Remaining uint16
	Reset     uint16
}

func (err *RateLimitError) Error() string {
	return fmt.Sprintf("%d of %d remaining requests. Next reset at %v", err.Remaining, err.Limit, err.Reset)
}

type parameters struct {
	Timezone string
	Rows     uint
	Format   string
	Staged   bool
}

type result struct {
	Nhits      uint
	Parameters parameters
	Datasets   []Dataset
}

type Dataset struct {
	Datasetid string
	Metas     struct {
		Publisher         string
		Domain            string
		RecordsCount      uint
		Title             string
		MetadataProcessed string
		DataProcessed     string
	}
	HasRecords bool
	fields     []interface{}
}

func DatasetCatalog(parameters *url.Values) []Dataset {

	result, err := fetchData("/api/datasets/1.0/search/", parameters)

	if err != nil {
		panic(err)
	}

	return result.Datasets

}

func LookupDataset() {

}

func DownloadRecords(filepath string, dataset Dataset) {

	filters := &url.Values{}
	filters.Set("dataset", dataset.Datasetid)
	DownloadFile(filepath, "/api/records/1.0/download", filters)
}

func LookupRecord() {

}

func fetchData(endpoint string, parameters *url.Values) (result, error) {

	// Get the data
	var result result

	resp, err := doRequest(endpoint, parameters)
	defer resp.Body.Close()

	if err != nil {
		return result, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
		return result, err
	}

	return result, nil
}

func DownloadFile(filepath string, endpoint string, parameters *url.Values) error {
	log.Println("Download file %s", filepath)

	// Get the data
	resp, err := doRequest(endpoint, parameters)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

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

func doRequest(endpoint string, parameters *url.Values) (*http.Response, error) {
	var bmurl = buildURL(endpoint, parameters)
	log.Println(bmurl)
	resp, err := http.Get(bmurl.String())

	if err != nil {
		return resp, err
	}

	if resp.StatusCode == 400 {
		resp.Body.Close()

		return resp, &RateLimitError{
			extractLimit(resp.Header, "X-Ratelimit-Limit"),
			extractLimit(resp.Header, "X-Ratelimit-Remaining"),
			extractLimit(resp.Header, "X-Ratelimit-Reset"),
		}
	}

	return resp, nil
}

func buildURL(endpoint string, parameters *url.Values) url.URL {
	return url.URL{
		Scheme:   "https",
		Host:     host,
		Path:     endpoint,
		RawQuery: parameters.Encode(),
	}
}

func extractLimit(headers map[string][]string, key string) uint16 {
	limit, err := strconv.ParseInt(headers[http.CanonicalHeaderKey(key)][0], 10, 16)

	if err != nil {
		panic(err)
	}

	return uint16(limit)
}
