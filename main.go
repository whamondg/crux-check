package main

import (
	"fmt"
	"os"
)

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

func main() {
	args := ReadArgs()

	var cruxRecord CruxRecord
	var err error

	cruxRecord, err = GetCruxData(args.apiKey, args.url, args.verbose)
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
