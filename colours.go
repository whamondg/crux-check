package main

import "fmt"

var green = "\033[32m"
var red = "\033[31m"
var grey = "\033[37m"

//DefaultColour for terminal output
var DefaultColour = grey

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
