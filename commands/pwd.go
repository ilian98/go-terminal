package commands

import (
	"fmt"
)

// Pwd is a structure for pwd command, implementing ExecuteCommand interface
type Pwd struct {
	path string
}

// GetName is a getter for command name
func (p *Pwd) GetName() string {
	return "pwd"
}

// GetPath is a getter for path
func (p *Pwd) GetPath() string {
	return p.path
}

// Clone is a method for cloning pwd command
func (p *Pwd) Clone() ExecuteCommand {
	clone := *p
	return &clone
}

// Execute is go implementation of pwd command
func (p *Pwd) Execute(cp CommandProperties) error {
	p.path = cp.Path
	inputFile, outputFile := cp.InputFile, cp.OutputFile
	defer closeInputOutputFiles(inputFile, outputFile)

	n, err := outputFile.WriteString(p.path)
	if err != nil {
		return err
	}
	if n != len(p.path) {
		return fmt.Errorf("Wrote only %d characters of current path: %s", n, p.path)
	}

	return nil
}
