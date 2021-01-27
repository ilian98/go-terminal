package commands

import (
	"fmt"
	"os"
)

// CommandProperties is used for storing the properties of command that will be executed
type CommandProperties struct {
	Path      string //  Path is used for storing the current path in the terminal for this command
	Arguments []string
	Options   []string
	Input     string // Empty Input would mean that we will use stdin for the command, otherwise it would be the name of the input file
	Output    string // Analogous to Input
}

// ExecuteCommand is interface for executing commands
type ExecuteCommand interface {
	GetName() string
	GetPath() string
	Execute(cp CommandProperties) error
	Clone() ExecuteCommand
}

func (cp *CommandProperties) pathFile(fileName string) string {
	return cp.Path + string(os.PathSeparator) + fileName
}

func (cp *CommandProperties) openInputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdin, nil
	}
	file, err := os.Open(cp.pathFile(fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("File for reading the input with name %s does not exist", fileName)
		}
		return nil, err
	}
	return file, nil
}

func (cp *CommandProperties) openOutputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdout, nil
	}
	file, err := os.OpenFile(cp.pathFile(fileName), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (cp *CommandProperties) openInputOutputFiles() (*os.File, *os.File, error) {
	inputFile, err := cp.openInputFile(cp.Input)
	if err != nil {
		return nil, nil, err
	}

	outputFile, err := cp.openOutputFile(cp.Output)
	if err != nil {
		return nil, nil, err
	}

	return inputFile, outputFile, nil
}

func (cp *CommandProperties) closeInputOutputFiles(inputFile *os.File, outputFile *os.File) {
	if cp.Input != "" {
		inputFile.Close()
	}
	if cp.Output != "" {
		outputFile.Close()
	}
}
