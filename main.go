package main

import (
	"encoding/json"
	"fmt"
	"os"
)

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
	}

	return "Fail"
}

func assessMetric(name string, metric Metric) MetricAssessment {
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
		assessMetric("CLS", cruxRecord.Record.Metrics.CLS),
		assessMetric("FID", cruxRecord.Record.Metrics.FID),
		assessMetric("LCP", cruxRecord.Record.Metrics.LCP),
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
