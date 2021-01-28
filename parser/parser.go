package parser

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Command is used for storing the properties of the inputted command after parsing
type Command struct {
	Name      string
	Options   []string
	Arguments []string
	Input     string // Empty Input would mean that we will use stdin for the command, otherwise it would be the name of the input file
	Output    string // Analogous to Input
	BgRun     bool
}

var (
	// ErrEmptyCommand indicates that the parsed command when trimmed is empty
	ErrEmptyCommand = errors.New("empty command")
)

// replaceEnclose replaces only those target characters with value that are enclosed in quotest
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

func parseCommandText(commandText string) (*Command, error) {
	commandText = strings.TrimSpace(commandText)

	commandText = replaceEnclosed(commandText, '\t', 0) // replace '\t' characters in probably arguments names with '\0' for save remove
	commandText = strings.Replace(commandText, "\t", " ", -1)
	commandText = replaceEnclosed(commandText, 0, '\t')

	commandText = replaceEnclosed(commandText, ' ', 0) // replace ' ' characters in probably arguments names with '\0' for save Split
	words := strings.Split(commandText, " ")
	// restoration of the ' ' characters
	for i, word := range words {
		words[i] = replaceEnclosed(word, 0, ' ')
	}

	if len(words) == 1 && len(words[0]) == 0 {
		return nil, ErrEmptyCommand
	}

	var c Command
	c.Name = words[0]
	c.BgRun = false

	removeQuotes := func(word string) string {
		if len(word) > 2 && word[0] == '"' && word[len(word)-1] == '"' {
			return word[1 : len(word)-1]
		}
		return word
	}
	for ind, word := range words[1:] {
		if len(word) == 0 {
			continue
		}
		if word == "&" {
			c.BgRun = true
			continue
		}

		if len(word) > 1 && word[0] == '-' {
			c.Options = append(c.Options, word[1:])
		} else if c.Input == "" && word[0] == '<' {
			// First argument with '<' will be considered for input, others will be counted as arguments
			if len(word) == 1 && (1+ind)+1 < len(words) {
				// This means file name is next argument
				c.Input = removeQuotes(words[(1+ind)+1])
				words[(1+ind)+1] = "" // This way we will skip it in next iteration
			} else {
				c.Input = removeQuotes(word[1:])
			}
		} else if c.Output == "" && word[0] == '>' {
			// First argument with '>' will be considered for output, others will be counted as arguments
			if len(word) == 1 && (1+ind)+1 < len(words) {
				// This means file name is next argument
				c.Output = removeQuotes(words[(1+ind)+1])
				words[(1+ind)+1] = "" // This way we will skip it in next iteration
			} else {
				c.Output = removeQuotes(word[1:])
			}
		} else {
			c.Arguments = append(c.Arguments, removeQuotes(word))
		}
	}
	return &c, nil
}

// Parse parses the string parameter text which should be an inputted command
func Parse(text string) ([]Command, error) {
	if runtime.GOOS == "windows" {
		text = strings.TrimRight(text, "\r\n")
	} else {
		text = strings.TrimRight(text, "\n")
	}

	var parsedCommand []Command

	text = replaceEnclosed(text, '|', 0) // replace '|' characters in probably arguments names with '\0' for save Split
	commandsText := strings.Split(text, "|")
	// restoration of the '|' characters
	for i, commandText := range commandsText {
		commandsText[i] = replaceEnclosed(commandText, 0, '|')
	}

	for _, commandText := range commandsText {
		command, err := parseCommandText(commandText)
		if err != nil {
			return nil, fmt.Errorf("Error when parsing command %s: %w", commandText, err)
		}
		parsedCommand = append(parsedCommand, *command)
	}
	return parsedCommand, nil
}
