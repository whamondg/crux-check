package main

import (
	"fmt"
	"os"
)

func analyseURL(url string, config Config) {
	cruxRecord, err := GetCruxData(config.apiKey, url, config.verbose)
	if err != nil {
		fmt.Printf("Failed to get CrUX data: %s\n", err)
		os.Exit(1)
	}

	metricResults := AssessCoreWebVitals(cruxRecord)

	fmt.Println()
	fmt.Printf("\033[1m%s\033[0m\n", cruxRecord.Record.Key.URL)
	fmt.Printf(Colourise("%10s:", ""), "Metric")
	fmt.Printf(Colourise("%6v", ""), "P75")
	fmt.Printf(Colourise("%11v", ""), "Threshold")
	fmt.Printf(Colourise("%8s\n", ""), "Status")

	for _, val := range metricResults {
		var scoreColour = FailColour
		if val.score == "Pass" {
			scoreColour = PassColour
		}

		fmt.Printf(Colourise("%10s:", ""), val.name)
		fmt.Printf(Colourise("%6v", ""), val.p75)
		fmt.Printf(Colourise("%10v", ""), val.threshold)
		fmt.Printf(Colourise("%8s\n", scoreColour), val.score)
	}
}

func main() {
	config := ReadConfig()
	for _, url := range config.urls {
		analyseURL(url, config)
	}
}
