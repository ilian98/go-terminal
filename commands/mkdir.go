package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	// ErrMkdirNoArgs indicates that there were no arguments passed to mkdir command
	ErrMkdirNoArgs = errors.New("At least one argument is needed")
	// ErrMkdirExists indicates that argument is name of file
	ErrMkdirExists = errors.New("exists")
)

// Mkdir is a structure for mkdir command, implementing ExecuteCommand interface
type Mkdir struct {
	path          string
	stopExecution chan struct{}
}

// GetName is a getter for command name
func (m *Mkdir) GetName() string {
	return "mkdir"
}

// GetPath is a getter for path
func (m *Mkdir) GetPath() string {
	return m.path
}

// Clone is a method for cloning mkdir command
func (m *Mkdir) Clone() ExecuteCommand {
	clone := *m
	return &clone
}

// InitChannel is a method for initializing stopExecution channel
func (m *Mkdir) InitChannel() {
	m.stopExecution = make(chan struct{}, 1)
}

// StopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (m *Mkdir) StopSignal() {
	m.stopExecution <- struct{}{}
}

// IsStopSignal is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (m *Mkdir) IsStopSignal() bool {
	select {
	case <-m.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of mkdir command
func (m *Mkdir) Execute(cp CommandProperties) error {
	m.path = cp.Path

	if len(cp.Arguments) == 0 {
		return ErrRmNoArgs
	}

	var errStrings []string
	for _, argument := range cp.Arguments {
		fullName := FullFileName(m.path, argument)
		_, err := os.Stat(fullName)

		if os.IsNotExist(err) {
			if err := os.Mkdir(fullName, 0666); err != nil {
				errStrings = append(errStrings, err.Error())
			}
		} else if err != nil {
			errStrings = append(errStrings, err.Error())
		} else {
			errStrings = append(errStrings, fmt.Errorf("%s - %w", fullName, ErrMkdirExists).Error())
		}
	}

	if len(errStrings) > 0 {
		return errors.New(strings.Join(errStrings, "\n"))
	}
	return nil
}
