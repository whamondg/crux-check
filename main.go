package main

import (
	"encoding/json"
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

/*
MetricAssessment contains an evaluation of a CrUX measurment scored against its target threshold
*/
type MetricAssessment struct {
	name      string
	valid     bool
	p75       string
	threshold string
	score     string
}

func validMetric(metric Metric) bool {
	return len(metric.Histogram) > 0
}

func metricScore(p75 json.Number, threshold json.Number) string {
	p75Num, err := p75.Float64()
	if err != nil {
		return "ERROR - Failed to convert p75 value to number"
	}

	thresholdNum, err := threshold.Float64()
	if err != nil {
		return "ERROR - Failed to convert threshold value to number"
	}

	if p75Num < thresholdNum {
		return "Pass"
	} else {
		return "Fail"
	}
}

func createMetricAssesment(name string, metric Metric) MetricAssessment {
	var p75, threshold, score string
	if validMetric(metric) {
		p75 = string(metric.Percentiles.P75)
		threshold = string(metric.Histogram[0].End)
		score = metricScore(metric.Percentiles.P75, metric.Histogram[0].End)
	}
	return MetricAssessment{p75: p75, threshold: threshold, name: name, score: score}
}

func assessCoreWebVitals(cruxRecord CruxRecord) []MetricAssessment {
	return []MetricAssessment{
		createMetricAssesment("CLS", cruxRecord.Record.Metrics.CLS),
		createMetricAssesment("FID", cruxRecord.Record.Metrics.FID),
		createMetricAssesment("LCP", cruxRecord.Record.Metrics.LCP),
	}
}

func main() {
	args := ReadArgs()

	cruxRecord, err := GetCruxData(args.apiKey, args.url, args.verbose)
	if err != nil {
		fmt.Printf("Failed to get CrUX data: %s\n", err)
		os.Exit(1)
	}

	metricResults := assessCoreWebVitals(cruxRecord)

	fmt.Println()
	fmt.Printf("%10s: %5v %10v %7s\n", "Metric", "P75", "Threshold", "Status")
	for _, val := range metricResults {
		fmt.Printf("%10s: %5v %10v %7s\n", val.name, val.p75, val.threshold, val.score)
	}
}
