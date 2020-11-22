package main

import "fmt"

var green = "\033[32m"
var red = "\033[31m"
var white = "\033[97m"

//DefaultColour for terminal output
var DefaultColour = white

//PassColour indicates a good result
var PassColour = green

//FailColour indicates a bad result
var FailColour = red

//Colourise creates a format string that will apply colour to the output
func Colourise(format string, colour string) string {
	if colour == "" {
		colour = DefaultColour
	}
	return fmt.Sprintf("%s%s%s", colour, format, DefaultColour)
}
