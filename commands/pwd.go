package commands

import (
	"fmt"
)

// Pwd is method of structure Execute command that executes go implementation of pwd command
func (e *ExecuteCommand) Pwd() error {
	inputFile, outputFile, err := e.openInputOutputFiles(e.Input, e.Output)
	defer e.closeInputOutputFiles(e.Input, inputFile, e.Output, outputFile)
	if err != nil {
		return err
	}

	n, err := outputFile.WriteString(e.Path)
	if err != nil {
		return err
	}
	if n != len(e.Path) {
		return fmt.Errorf("Wrote only %d characters of current path: %s", n, e.Path)
	}

	return nil
}
