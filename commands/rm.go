package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	// ErrRmNoArgs indicates that there were no arguments passed to rm command
	ErrRmNoArgs = errors.New("At least one argument is needed")
	// ErrRmIsFile indicates that argument is name of file
	ErrRmIsFile = errors.New("is a file")
	// ErrRmIsDir indicates that argument is name of directory
	ErrRmIsDir = errors.New("is a directory")
	// ErrRmInvalidName indicates that argument is not a valid name in file system
	ErrRmInvalidName = errors.New("is not a valid name in the file system")
)

// Rm is a structure for rm command, implementing ExecuteCommand interface
type Rm struct {
	path          string
	stopExecution chan struct{}
}

// GetName is a getter for command name
func (r *Rm) GetName() string {
	return "rm"
}

// GetPath is a getter for path
func (r *Rm) GetPath() string {
	return r.path
}

// Clone is a method for cloning pwd command
func (r *Rm) Clone() ExecuteCommand {
	clone := *r
	return &clone
}

// InitStopCatching is a method for initializing stopExecution channel
func (r *Rm) InitStopCatching() {
	r.stopExecution = make(chan struct{}, 1)
}

// StopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (r *Rm) StopSignal() {
	r.stopExecution <- struct{}{}
}

// IsStopSignal is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (r *Rm) IsStopSignal() bool {
	select {
	case <-r.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of rm command
func (r *Rm) Execute(cp CommandProperties) error {
	r.path = cp.Path

	if len(cp.Arguments) == 0 {
		return ErrRmNoArgs
	}

	recursiveOption := false
	for _, option := range cp.Options {
		if option == "r" {
			recursiveOption = true
		}
	}

	var errStrings []string
	for _, argument := range cp.Arguments {
		fullName := FullFileName(r.path, argument)
		stat, err := os.Stat(fullName)

		if os.IsNotExist(err) {
			errStrings = append(errStrings, fmt.Errorf("%s %w", fullName, ErrRmInvalidName).Error())
			continue
		}
		if err != nil {
			errStrings = append(errStrings, err.Error())
			continue
		}

		if recursiveOption == true {
			if stat.IsDir() {
				if err := os.RemoveAll(fullName); err != nil {
					errStrings = append(errStrings, err.Error())
				}
			} else {
				errStrings = append(errStrings, fmt.Errorf("%s %w", fullName, ErrRmIsFile).Error())
			}
		} else {
			if !stat.IsDir() {
				if err := os.Remove(fullName); err != nil {
					errStrings = append(errStrings, err.Error())
				}
			} else {
				errStrings = append(errStrings, fmt.Errorf("%s %w", fullName, ErrRmIsDir).Error())
			}
		}
	}

	if len(errStrings) > 0 {
		return errors.New(strings.Join(errStrings, "\n"))
	}
	return nil
}
