package commands

import (
	"os"
)

// CommandProperties is used for storing the properties of command that will be executed
type CommandProperties struct {
	Path       string //  Path is used for storing the current path in the terminal for this command
	Arguments  []string
	Options    []string
	InputFile  *os.File // InputFile is used for reading the input it could be stdin
	OutputFile *os.File // OutputFile is used for reading the input it could be stdin
}

// ExecuteCommand is interface for executing commands
type ExecuteCommand interface {
	GetName() string
	GetPath() string
	Execute(cp CommandProperties) error
	Clone() ExecuteCommand
}

func closeInputOutputFiles(inputFile *os.File, outputFile *os.File) {
	if inputFile != os.Stdin {
		inputFile.Close()
	}
	if outputFile != os.Stdout {
		outputFile.Close()
	}
}
