package main

import (
	"flag"
	"fmt"
	"os"
)

/*
Config is a struct for storing runtime configuration
*/
type Config struct {
	apiKey  string
	urls    []string
	verbose bool
}

/*
ReadConfig parses command line arguments and looks up required environment variables.
*/
func ReadConfig() Config {
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

	return Config{apiKey: apiKey, urls: []string{*url}, verbose: *verbose}
}
