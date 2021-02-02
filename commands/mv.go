package commands

import (
	"errors"
	"fmt"
	"os"
)

var (
	// ErrMvTwoArgs indicates that the argument count is different from two
	ErrMvTwoArgs = errors.New("Two arguments are needed")
	// ErrMvInvalidName indicates that source argument is not a valid name in file system
	ErrMvInvalidName = errors.New("is not a valid name in the file system")
	// ErrMvIsDir indicates that source argument is name of directory
	ErrMvIsDir = errors.New("is a directory")
	// ErrMvSame indicates that source and target are the same file
	ErrMvSame = errors.New("Source and target file are the same")
)

// Mv is a structure for mv command, implementing ExecuteCommand interface
type Mv struct {
	path          string
	stopExecution chan struct{}
}

// GetName is a getter for command name
func (m *Mv) GetName() string {
	return "mv"
}

// GetPath is a getter for path
func (m *Mv) GetPath() string {
	return m.path
}

// Clone is a method for cloning mv command
func (m *Mv) Clone() ExecuteCommand {
	clone := *m
	return &clone
}

// InitStopSignalCatching is a method for initializing stopExecution channel
func (m *Mv) InitStopSignalCatching() {
	m.stopExecution = make(chan struct{}, 1)
}

// SendStopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (m *Mv) SendStopSignal() {
	m.stopExecution <- struct{}{}
}

// IsStopSignalReceived is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (m *Mv) IsStopSignalReceived() bool {
	select {
	case <-m.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of mv command
func (m *Mv) Execute(cp CommandProperties) error {
	m.path = cp.Path

	if len(cp.Arguments) != 2 {
		return ErrMvTwoArgs
	}

	source := FullFileName(m.path, cp.Arguments[0])
	dest := FullFileName(m.path, cp.Arguments[1])
	if source == dest {
		return ErrMvSame
	}

	stat, err := os.Stat(source)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s - %w", source, ErrMvInvalidName)
	} else if err != nil {
		return err
	} else if stat.IsDir() {
		return fmt.Errorf("%s - %w", source, ErrMvIsDir)
	}

	if err := os.Rename(source, dest); err != nil {
		return err
	}

	return nil
}
