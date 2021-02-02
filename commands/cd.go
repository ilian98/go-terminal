package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	// ErrCdTooManyArgs indicates that the cd command has more than 1 argument
	ErrCdTooManyArgs = errors.New("Too many arguments")
	// ErrCdPathLeadsToFile indicates that the path in cd command is to file
	ErrCdPathLeadsToFile = errors.New("path leads to file, not directory")
	// ErrCdPathNotExist indicates that the path in cd command does not exist
	ErrCdPathNotExist = errors.New("path does not exist")
)

// Cd is a structure for cd command, implementing ExecuteCommand interface
type Cd struct {
	path          string
	stopExecution chan struct{}
}

// GetName is a getter for command name
func (c *Cd) GetName() string {
	return "cd"
}

// GetPath is a getter for path
func (c *Cd) GetPath() string {
	return c.path
}

// Clone is a method for cloning cd command
func (c *Cd) Clone() ExecuteCommand {
	clone := *c
	return &clone
}

// InitStopSignalCatching is a method for initializing stopExecution channel
func (c *Cd) InitStopSignalCatching() {
	c.stopExecution = make(chan struct{}, 1)
}

// SendStopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (c *Cd) SendStopSignal() {
	c.stopExecution <- struct{}{}
}

// IsStopSignalReceived is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (c *Cd) IsStopSignalReceived() bool {
	select {
	case <-c.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of cd command
func (c *Cd) Execute(cp CommandProperties) error {
	c.path = cp.Path

	if len(cp.Arguments) == 0 {
		c.path = getRootPath(c.path)
		return nil
	}
	if len(cp.Arguments) > 1 {
		return ErrCdTooManyArgs
	}

	path := cp.Arguments[0]
	if len(path) == 0 {
		c.path = getRootPath(c.path)
		return nil
	}

	tryPath := FullFileName(c.path, path)

	stat, err := os.Stat(tryPath) // we use stat to check if the path is valid and leading to directory
	if err == nil && stat.IsDir() {
		p, err := filepath.Abs(tryPath)
		if err != nil {
			return err
		}
		c.path = p
		return nil
	} else if err == nil {
		return fmt.Errorf("%s - %w", tryPath, ErrCdPathLeadsToFile)
	} else if os.IsNotExist(err) {
		return fmt.Errorf("%s - %w", tryPath, ErrCdPathNotExist)
	}
	return err
}
