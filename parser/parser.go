package parser

import (
	"runtime"
	"strings"
)

// Command is used for storing the properties of the inputted command after parsing
type Command struct {
	Name      string
	Options   []string
	Arguments []string
}

// Parse parses the string parameter text which should be an inputted command
func Parse(text string) []Command {
	if runtime.GOOS == "windows" {
		text = strings.TrimRight(text, "\r\n")
	} else {
		text = strings.TrimRight(text, "\n")
	}

	var c Command
	c.Name = text
	parsedCommand := make([]Command, 1)
	parsedCommand[0] = c
	return parsedCommand
}
