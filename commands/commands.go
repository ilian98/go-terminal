package commands

import (
	"errors"
	"fmt"
	"os"
)

// Path is a global variable for storing the current path in the terminal
var Path string

func openInputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdin, nil
	}
	file, err := os.Open(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("File for storing the input with name %s does not exist", fileName)
		}
		return nil, err
	}
	return file, nil
}

func openOutputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdout, nil
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func openInputOutputFiles(inputName string, outputName string) (*os.File, *os.File, error) {
	inputFile, err := openInputFile(inputName)
	if err != nil {
		return nil, nil, err
	}

	outputFile, err := openOutputFile(outputName)
	if err != nil {
		return nil, nil, err
	}

	return inputFile, outputFile, nil
}

func closeInputOutputFiles(input string, inputFile *os.File, output string, outputFile *os.File) {
	if input != "" {
		inputFile.Close()
	}
	if output != "" {
		outputFile.Close()
	}
}
