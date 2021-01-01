package parser

import "fmt"

// command is used for storing the properties of the inputted command after parsing
type command struct {
	name      string
	options   []string
	arguments []string
}

func parse(text string) []command {
	fmt.Println(text)
	var c command
	c.name = text
	parsedCommand := make([]command, 1)
	parsedCommand[0] = c
	return parsedCommand
}
