package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

var (
	// ErrCatFileNotExist indicates that one of the arguments is invalid file name
	ErrCatFileNotExist = errors.New("file does not exist")
)

// Cat is a structure for cat command, implementing ExecuteCommand interface
type Cat struct {
	path string
}

// GetName is a getter for command name
func (c *Cat) GetName() string {
	return "cat"
}

// GetPath is a getter for path
func (c *Cat) GetPath() string {
	return c.path
}

// Clone is a method for cloning cat command
func (c *Cat) Clone() ExecuteCommand {
	clone := *c
	return &clone
}

// Execute is go implementation of cat command
func (c *Cat) Execute(cp CommandProperties) error {
	c.path = cp.Path
	inputFile, outputFile := cp.InputFile, cp.OutputFile
	defer closeInputOutputFiles(inputFile, outputFile)

	outputFileData := func(file *os.File) error {
		for {
			buffer := make([]byte, 1<<4)
			n, err := file.Read(buffer)
			if n == 0 {
				break
			}
			text := string(buffer)
			outputFile.WriteString(strings.TrimRight(text, "\u0000"))
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
		}
		return nil
	}

	if len(cp.Arguments) == 0 {
		err := outputFileData(inputFile)
		if err != nil {
			return err
		}

		// clean newline after EOF if the reading was from stdin
		var bufferNewLine []byte
		if runtime.GOOS == "windows" {
			bufferNewLine = make([]byte, 2)
		} else {
			bufferNewLine = make([]byte, 1)
		}
		if _, err := inputFile.Read(bufferNewLine); err != nil && err != io.EOF {
			return err
		}
		return nil
	}

	var errStrings []string
	for _, argument := range cp.Arguments {
		file, err := os.Open(c.path + string(os.PathSeparator) + argument)
		if os.IsNotExist(err) {
			errStrings = append(errStrings, fmt.Errorf("%s - %w", argument, ErrCatFileNotExist).Error())
		} else if err != nil {
			errStrings = append(errStrings, fmt.Errorf("%s - %w", argument, err).Error())
		} else {
			err := outputFileData(file)
			file.Close()
			if err != nil {
				errStrings = append(errStrings, fmt.Errorf("%s - %w", argument, err).Error())
			}
		}
	}

	if len(errStrings) == 0 {
		return nil
	}
	return errors.New(strings.Join(errStrings, "\n"))
}
