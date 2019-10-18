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

//RateLimitError : Api Quotas https://help.opendatasoft.com/apis/ods-search-v1/#quotas
type RateLimitError struct {
	Limit     uint16
	Remaining uint16
	Reset     uint16
}

func (err *RateLimitError) Error() string {
	return fmt.Sprintf("%d of %d remaining requests. Next reset at %v", err.Remaining, err.Limit, err.Reset)
}

//Parameters result parameters
type Parameters struct {
	Timezone string
	Rows     uint
	Format   string
	Staged   bool
}

//FetchDatasetsCatalog Fetch a catalog of datasets
func FetchDatasetsCatalog(parameters *url.Values, result interface{}) {

	err := fetchData("/api/datasets/1.0/search/", parameters, result)

	if err != nil {
		panic(err)
	}
}

//DownloadDataset download a dataset into a dest file
func DownloadDataset(dataset Dataset, dest string) error {
	filters := &url.Values{}
	filters.Set("dataset", dataset.Datasetid)
	return downloadFile("/api/records/1.0/download", filters, dest)
}

func fetchData(endpoint string, parameters *url.Values, result interface{}) error {

	resp, err := doRequest(endpoint, parameters)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}

	return nil
}

func downloadFile(endpoint string, parameters *url.Values, dest string) error {
	log.Printf("Download inprogress : %v", endpoint)

	// Get the data
	resp, err := doRequest(endpoint, parameters)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	log.Printf("Download completed : %v", dest)

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
