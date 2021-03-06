package commands

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrFindNoArgs indicates that there were no arguments passed to find command
	ErrFindNoArgs = errors.New("At least one argument is needed")
	// ErrFindFound indicates that the file was found in walk
	ErrFindFound = errors.New("found")
)

// Find is a structure for pwd command, implementing ExecuteCommand interface
type Find struct {
	path          string
	stopExecution chan struct{}
}

// GetName is a getter for command name
func (f *Find) GetName() string {
	return "find"
}

// GetPath is a getter for path
func (f *Find) GetPath() string {
	return f.path
}

// Clone is a method for cloning find command
func (f *Find) Clone() ExecuteCommand {
	clone := *f
	return &clone
}

// InitStopSignalCatching is a method for initializing stopExecution channel
func (f *Find) InitStopSignalCatching() {
	f.stopExecution = make(chan struct{}, 1)
}

// SendStopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (f *Find) SendStopSignal() {
	f.stopExecution <- struct{}{}
}

// IsStopSignalReceived is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (f *Find) IsStopSignalReceived() bool {
	select {
	case <-f.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of find command
func (f *Find) Execute(cp CommandProperties) error {
	f.path = cp.Path
	_, outputFile := cp.InputFile, cp.OutputFile

	if len(cp.Arguments) == 0 {
		return ErrFindNoArgs
	}

	var errStrings []string
	returnFunc := func() error {
		if len(errStrings) == 0 {
			return nil
		}
		return errors.New(strings.Join(errStrings, "\n"))
	}

	for _, argument := range cp.Arguments {
		if f.IsStopSignalReceived() == true {
			return returnFunc()
		}

		err := filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
			if f.IsStopSignalReceived() == true {
				return ErrStoppedExec
			}
			if err != nil {
				return err
			}
			s := strings.Split(path, string(os.PathSeparator)) // we split the path so we can only check the short name
			if len(s) > 0 && s[len(s)-1] == argument {
				if err := checkWrite(f, outputFile, argument+" found - "+path+"\n"); err != nil {
					return err
				}
				return ErrFindFound
			}
			return nil
		})
		if err == nil {
			if err := checkWrite(f, outputFile, argument+" not found\n"); err != nil {
				return err
			}
		} else if err == ErrStoppedExec {
			return returnFunc()
		} else if err != ErrFindFound {
			errStrings = append(errStrings, err.Error())
		}
	}

	return returnFunc()
}
