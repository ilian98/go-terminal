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
	path          string
	stopExecution chan struct{}
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

// InitStopSignalCatching is a method for initializing stopExecution channel
func (c *Cat) InitStopSignalCatching() {
	c.stopExecution = make(chan struct{}, 1)
}

// SendStopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (c *Cat) SendStopSignal() {
	c.stopExecution <- struct{}{}
}

// IsStopSignalReceived is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (c *Cat) IsStopSignalReceived() bool {
	select {
	case <-c.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of cat command
func (c *Cat) Execute(cp CommandProperties) error {
	c.path = cp.Path
	inputFile, outputFile := cp.InputFile, cp.OutputFile

	outputFileData := func(file *os.File) error {
		for {
			text, err := checkRead(c, file)
			if len(text) == 0 {
				break
			}
			if err != nil && err != io.EOF {
				return err
			}
			if err := checkWrite(c, outputFile, text); err != nil {
				return err
			}
			if err == io.EOF {
				break
			}
		}
		return nil
	}

	if len(cp.Arguments) == 0 { // when there are no arguments, cat command reads from inputFile
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

	var errStrings []string // in slice errStrings we collect all the errors
	for _, argument := range cp.Arguments {
		file, err := os.Open(c.path + string(os.PathSeparator) + argument)
		if os.IsNotExist(err) {
			errStrings = append(errStrings, fmt.Errorf("%s - %w", argument, ErrCatFileNotExist).Error())
		} else if err != nil {
			errStrings = append(errStrings, fmt.Errorf("%s - %w", argument, err).Error())
		} else {
			err := outputFileData(file)
			file.Close()
			if err == ErrStoppedExec {
				if len(errStrings) == 0 {
					return nil
				}
				return errors.New(strings.Join(errStrings, "\n"))
			}
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
