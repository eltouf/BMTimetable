//Package to fetch data from https://opendata.bordeaux-metropole.fr/api/datasets
package main

import (
	"BMTimetable/internal/client"
	"net/url"
)

func main() {
	filters := &url.Values{}
	filters.Set("rows", "50")
	filters.Set("refine.keyword", "saeiv")

	client.FetchDatasets(filters)
}
