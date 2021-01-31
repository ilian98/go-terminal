package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	// ErrCpTwoArgs indicates that the argument count is different from two
	ErrCpTwoArgs = errors.New("Two arguments are needed")
	// ErrCpInvalidName indicates that source argument is not a valid name in file system
	ErrCpInvalidName = errors.New("is not a valid name in the file system")
	// ErrCpIsDir indicates that source argument is name of directory
	ErrCpIsDir = errors.New("is a directory")
	// ErrCpSame indicates that source and target are the same file
	ErrCpSame = errors.New("Source and target file are the same")
)

// Cp is a structure for cp command, implementing ExecuteCommand interface
type Cp struct {
	path          string
	stopExecution chan struct{}
}

// GetName is a getter for command name
func (c *Cp) GetName() string {
	return "cp"
}

// GetPath is a getter for path
func (c *Cp) GetPath() string {
	return c.path
}

// Clone is a method for cloning cp command
func (c *Cp) Clone() ExecuteCommand {
	clone := *c
	return &clone
}

// InitChannel is a method for initializing stopExecution channel
func (c *Cp) InitChannel() {
	c.stopExecution = make(chan struct{}, 1)
}

// StopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (c *Cp) StopSignal() {
	c.stopExecution <- struct{}{}
}

// IsStopSignal is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (c *Cp) IsStopSignal() bool {
	select {
	case <-c.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of cp command
func (c *Cp) Execute(cp CommandProperties) error {
	c.path = cp.Path

	if len(cp.Arguments) != 2 {
		return ErrCpTwoArgs
	}

	source := FullFileName(c.path, cp.Arguments[0])
	dest := FullFileName(c.path, cp.Arguments[1])
	if source == dest {
		return ErrCpSame
	}

	stat, err := os.Stat(source)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s - %w", source, ErrCpInvalidName)
	} else if err != nil {
		return err
	} else if stat.IsDir() {
		return fmt.Errorf("%s - %w", source, ErrCpIsDir)
	}
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	copy, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer copy.Close()

	_, err = io.Copy(copy, file)
	if err != nil {
		return err
	}
	return nil
}
