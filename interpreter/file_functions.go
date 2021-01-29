package interpreter

import (
	"fmt"
	"os"

	"github.com/ilian98/go-terminal/commands"
)

func (i *Interpreter) openInputFile(fileName string) (*os.File, error) {
	if fileName == "" {
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

func (i *Interpreter) openOutputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdout, nil
	}
	file, err := os.OpenFile(commands.FullFileName(i.Path, fileName), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

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

func closeInputOutputFiles(inputFile *os.File, outputFile *os.File) {
	if inputFile != os.Stdin {
		inputFile.Close()
	}
	if outputFile != os.Stdout {
		outputFile.Close()
	}
}
