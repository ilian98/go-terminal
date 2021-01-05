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

// Equal is a method for checking if two Commands are equal
func (c1 *Command) Equal(c2 Command) bool {
	if c1.Name != c2.Name {
		return false
	}
	if len(c1.Options) != len(c2.Options) {
		return false
	}
	for ind := range c1.Options {
		if c1.Options[ind] != c2.Options[ind] {
			return false
		}
	}
	if len(c1.Arguments) != len(c2.Arguments) {
		return false
	}
	for ind := range c1.Arguments {
		if c1.Arguments[ind] != c2.Arguments[ind] {
			return false
		}
	}
	return true
}

// NotEqual is method for checking if two Commands are different
func (c1 *Command) NotEqual(c2 Command) bool {
	return !c1.Equal(c2)
}

// replaceEnclose replaces only those target characters with value that are enclosed in "" ""
func replaceEnclosed(text string, target byte, value byte) string {
	byteText := []byte(text)
	quotes := 0
	for ind, char := range byteText {
		if char == target && quotes%2 == 1 {
			byteText[ind] = value
		}
		if char == '"' {
			quotes++
		}
	}
	return string(byteText)
}

func parseCommandText(commandText string) Command {
	commandText = strings.Trim(commandText, " \t")
	var c Command

	commandText = replaceEnclosed(commandText, ' ', 0) // replace ' ' characters in probably arguments names with '\0' for save Split
	words := strings.Split(commandText, " ")
	// restoration of the ' ' characters
	for i, word := range words {
		words[i] = replaceEnclosed(word, 0, ' ')
	}

	if len(words) == 0 {
		// Error here
	}
	c.Name = words[0]
	c.Options = make([]string, 0)
	c.Arguments = make([]string, 0)
	for _, word := range words[1:] {
		word = strings.Trim(word, "\t")
		if len(word) == 0 {
			continue
		}
		if word[0] == '-' {
			c.Options = append(c.Options, word[1:])
		} else {
			if len(word) > 2 && word[0] == '"' && word[len(word)-1] == '"' {
				c.Arguments = append(c.Arguments, word[1:len(word)-1])
			} else {
				c.Arguments = append(c.Arguments, word)
			}
		}
	}
	return c
}

// Parse parses the string parameter text which should be an inputted command
func Parse(text string) []Command {
	if runtime.GOOS == "windows" {
		text = strings.TrimRight(text, "\r\n")
	} else {
		text = strings.TrimRight(text, "\n")
	}

	parsedCommand := make([]Command, 0)

	text = replaceEnclosed(text, '|', 0) // replace '|' characters in probably arguments names with '\0' for save Split
	commandsText := strings.Split(text, "|")
	// restoration of the '|' characters
	for i, commandText := range commandsText {
		commandsText[i] = replaceEnclosed(commandText, 0, '|')
	}

	for _, commandText := range commandsText {
		parsedCommand = append(parsedCommand, parseCommandText(commandText))
	}
	return parsedCommand
}
