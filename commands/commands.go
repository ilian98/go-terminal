package commands

import (
	"fmt"
	"os"
)

// ExecuteCommand is used for storing the properties and executing parsed command
type ExecuteCommand struct {
	Path      string //  Path is used for storing the current path in the terminal for this command
	Arguments []string
	Options   []string
	Input     string // Empty Input would mean that we will use stdin for the command, otherwise it would be the name of the input file
	Output    string // Analogous to Input
}

func (e *ExecuteCommand) pathFile(fileName string) string {
	return e.Path + string(os.PathSeparator) + fileName
}

func (e *ExecuteCommand) openInputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdin, nil
	}
	file, err := os.Open(e.pathFile(fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("File for storing the input with name %s does not exist", fileName)
		}
		return nil, err
	}
	return file, nil
}

func (e *ExecuteCommand) openOutputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdout, nil
	}
	file, err := os.OpenFile(e.pathFile(fileName), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (e *ExecuteCommand) openInputOutputFiles(inputName string, outputName string) (*os.File, *os.File, error) {
	inputFile, err := e.openInputFile(inputName)
	if err != nil {
		return nil, nil, err
	}

	outputFile, err := e.openOutputFile(outputName)
	if err != nil {
		return nil, nil, err
	}

	return inputFile, outputFile, nil
}

func (e *ExecuteCommand) closeInputOutputFiles(input string, inputFile *os.File, output string, outputFile *os.File) {
	if input != "" {
		inputFile.Close()
	}
	if output != "" {
		outputFile.Close()
	}
}
