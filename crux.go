package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type CruxRecord struct {
	Record struct {
		Metrics struct {
			CLS Metric `json:"cumulative_layout_shift"`
			FID Metric `json:"first_input_delay"`
			LCP Metric `json:"largest_contentful_paint"`
		} `json:"metrics"`
	} `json:"record"`
}

type Metric struct {
	Histogram   []Histogram `json:"histogram"`
	Percentiles Percentiles `json:"percentiles"`
}

type Histogram struct {
	Density json.Number `json:"density"`
	Start   json.Number `json:"start"`
	End     json.Number `json:"end,omitempty"`
}

type Percentiles struct {
	P75 json.Number `json:"p75"`
}

func getCruxData(apiKey string, target string, verbose bool) (CruxRecord, error) {
	cruxRecord := CruxRecord{}

	fmt.Printf("Checking CrUX data for %s\n", target)

	cruxAPI := fmt.Sprintf("https://chromeuxreport.googleapis.com/v1/records:queryRecord?key=%s", apiKey)
	requestBody, _ := json.Marshal(map[string]string{"url": target})
	response, err := http.Post(cruxAPI, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return cruxRecord, fmt.Errorf("The HTTP request failed with error %s", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		if verbose {
			fmt.Printf("Received CrUX data\n")
			fmt.Println(string(data))
		}

		jsonErr := json.Unmarshal(data, &cruxRecord)
		if jsonErr != nil {
			return cruxRecord, fmt.Errorf("JSON Unmarshalling failed with error %s", err)
		}
		return cruxRecord, nil
	}
}

func assessMetric(name string, metric Metric) {
	// fmt.Printf("%+v\n", metric)

	p75 := metric.Percentiles.P75
	threshold := metric.Histogram[0].End
	var result = "fail"
	// if (p75.Float64()) < (target.Float64()) {
	// 	result = "pass"
	// }

	fmt.Printf("%7s: %5v %10v %7s\n", name, p75, threshold, result)
}

func main() {
	apiKey := os.Getenv("CRUX_API_KEY")
	if len(apiKey) == 0 {
		fmt.Fprintf(os.Stderr, "No API key defined: set the CRUX_API_KEY environment variable\n")
		os.Exit(1)
	}

	url := flag.String("u", "", "A URL which will be used to find CrUX data")
	verbose := flag.Bool("v", false, "Enable verbose logging")
	required := []string{"u"}

	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })

	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "mandatory -%s flag is missing\n\n", req)
			flag.Usage()
			os.Exit(1)
		}
	}

	cruxRecord, err := getCruxData(apiKey, *url, *verbose)
	if err != nil {
		fmt.Printf("Failed to get CrUX data: %s\n", err)
	}

	if *verbose {
		fmt.Printf("Parsed JSON into Structs\n")
		fmt.Printf("%+v\n", cruxRecord)
	}

	fmt.Println()
	fmt.Printf("%7s: %5v %10v %7s\n", "Metric", "P75", "Threshold", "Status")
	assessMetric("CLS", cruxRecord.Record.Metrics.CLS)
	assessMetric("FID", cruxRecord.Record.Metrics.FID)
	assessMetric("LCP", cruxRecord.Record.Metrics.LCP)
}
