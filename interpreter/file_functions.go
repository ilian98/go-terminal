package interpreter

import (
	"fmt"
	"os"

	"github.com/ilian98/go-terminal/commands"
)

// openInputFile is a function that opens file for input and checks if the fileName is a relative path
func (i *Interpreter) openInputFile(fileName string) (*os.File, error) {
	if fileName == "" { // empty file name means that input file should be os.Stdin
		return os.Stdin, nil
	}
	file, err := os.Open(commands.FullFileName(i.Path, fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("File for reading the input with name %s does not exist", fileName)
		}
		return nil, err
	}
	return file, nil
}

// openOutputFile is a function that opens file for output and checks if the fileName is a relative path
func (i *Interpreter) openOutputFile(fileName string) (*os.File, error) {
	if fileName == "" { // empty file name means that output file should be os.Stdout
		return os.Stdout, nil
	}
	file, err := os.OpenFile(commands.FullFileName(i.Path, fileName), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// openInputOutpuFiles is a function that calls openInputFile and openOutputFile for opening files for input and output
func (i *Interpreter) openInputOutputFiles(input string, output string) (*os.File, *os.File, error) {
	inputFile, err := i.openInputFile(input)
	if err != nil {
		return nil, nil, err
	}

	outputFile, err := i.openOutputFile(output)
	if err != nil {
		return nil, nil, err
	}

	return inputFile, outputFile, nil
}

// closeInputOutputFiles is function for closing opened inputFile and outputFile if they are different from os.Stdin and os.Stdout respectively
func closeInputOutputFiles(inputFile *os.File, outputFile *os.File) {
	if inputFile != os.Stdin {
		inputFile.Close()
	}
	if outputFile != os.Stdout {
		outputFile.Close()
	}
}
