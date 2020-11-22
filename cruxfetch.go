package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
CruxRecord represents the overall CrUX api response
*/
type CruxRecord struct {
	Record struct {
		Metrics struct {
			CLS Metric `json:"cumulative_layout_shift"`
			FID Metric `json:"first_input_delay"`
			LCP Metric `json:"largest_contentful_paint"`
		} `json:"metrics"`
	} `json:"record"`
}

/*
Metric stores data about a specific CrUX measurement
*/
type Metric struct {
	Histogram   []Histogram `json:"histogram"`
	Percentiles Percentiles `json:"percentiles"`
}

/*
Histogram splits the Metric data into Good, Needs Improvement, and Poor
*/
type Histogram struct {
	Density json.Number `json:"density"`
	Start   json.Number `json:"start"`
	End     json.Number `json:"end,omitempty"`
}

/*
Percentiles provieds the Metric measurement for a specific percentile
*/
type Percentiles struct {
	P75 json.Number `json:"p75"`
}

/*
GetCruxData calls the CrUX api to retrieve user experience data
*/
func GetCruxData(apiKey string, target string, verbose bool) (CruxRecord, error) {
	cruxRecord := CruxRecord{}

	fmt.Printf("Checking CrUX data for %s\n", target)

	cruxAPI := fmt.Sprintf("https://chromeuxreport.googleapis.com/v1/records:queryRecord?key=%s", apiKey)
	requestBody, _ := json.Marshal(map[string]string{"url": target})
	response, err := http.Post(cruxAPI, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return cruxRecord, fmt.Errorf("The HTTP request failed with error %s", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode != 200 {
			return cruxRecord, fmt.Errorf("Received an error while fetching CrUX data\n %s", string(data))
		}

		if verbose {
			fmt.Printf("Received CrUX data\n")
			fmt.Println(string(data))
		}

		jsonErr := json.Unmarshal(data, &cruxRecord)
		if jsonErr != nil {
			return cruxRecord, fmt.Errorf("JSON Unmarshalling failed with error %s", err)
		}

		if verbose {
			fmt.Printf("Parsed JSON into Structs\n")
			fmt.Printf("%+v\n", cruxRecord)
		}

		return cruxRecord, nil
	}
}
