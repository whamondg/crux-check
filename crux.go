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

type CliArgs struct {
	apiKey  string
	url     string
	verbose bool
}

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

func assessMetric(metric Metric) (bool, error) {
	p75Num, err := metric.Percentiles.P75.Float64()
	if err != nil {
		return false, fmt.Errorf("Failed to convert p75 value of %s to float64: %s", metric.Percentiles.P75, err)
	}

	thresholdNum, err := metric.Histogram[0].End.Float64()
	if err != nil {
		return false, fmt.Errorf("Failed to convert threshold value of %s to float64: %s", metric.Histogram[0].End, err)
	}

	return p75Num < thresholdNum, nil
}

func readArgs() CliArgs {
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

	return CliArgs{apiKey: apiKey, url: *url, verbose: *verbose}
}

func main() {
	args := readArgs()

	var cruxRecord CruxRecord
	var err error

	cruxRecord, err = getCruxData(args.apiKey, args.url, args.verbose)
	if err != nil {
		fmt.Printf("Failed to get CrUX data: %s\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("%10s: %5v %10v %7s\n", "Metric", "P75", "Threshold", "Status")

	var status bool
	status, err = assessMetric(cruxRecord.Record.Metrics.CLS)
	if err != nil {
		fmt.Printf("Failed to asses CLS: %s\n", err)
		os.Exit(1)
	}

	printMetric(cruxRecord.Record.Metrics.CLS, status)

	status, err = assessMetric(cruxRecord.Record.Metrics.FID)
	if err != nil {
		fmt.Printf("Failed to asses FID: %s\n", err)
		os.Exit(1)
	}
	printMetric(cruxRecord.Record.Metrics.FID, status)

	status, err = assessMetric(cruxRecord.Record.Metrics.LCP)
	if err != nil {
		fmt.Printf("Failed to asses LCP: %s\n", err)
		os.Exit(1)
	}

	printMetric(cruxRecord.Record.Metrics.LCP, status)
}

func printMetric(metric Metric, status bool) {
	var result = "fail"
	if status {
		result = "pass"
	}

	fmt.Printf("%10s: %5v %10v %7s\n", "CLS", metric.Percentiles.P75, metric.Histogram[0].End, result)
}
