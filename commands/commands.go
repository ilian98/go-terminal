// Package commands defines the interface for commands, some helper functions and implements the commands.
package commands

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	// ErrStoppedExec indicates that the execution of command was stopped
	ErrStoppedExec = errors.New("execution was stopped")
)

// CommandProperties is used for storing the properties of command that will be executed
type CommandProperties struct {
	Path       string //  Path is used for storing the current path in the terminal for this command
	Arguments  []string
	Options    []string
	InputFile  *os.File // InputFile is used for reading the input it could be stdin
	OutputFile *os.File // OutputFile is used for reading the input it could be stdin
}

// Function for constructing CommandProperties object with only path, arguments and options
func newCp(Path string, Arguments []string, Options []string) CommandProperties {
	return CommandProperties{Path, Arguments, Options, os.Stdin, os.Stdout}
}

// ExecuteCommand is interface for executing commands
//
// The interface includes getters for name and path.
// The Clone method is important - it allows the command to run clean every time by cloning the initial state in interpreter package
//
// The method InitStopCatching should be used for initializing the catching of stop signals.
// The method StopSignal should be used outside (from package interpreter) to send stop signal.
// The method IsStopSignal should be used by command to check if stop signal is received.
type ExecuteCommand interface {
	GetName() string
	GetPath() string
	Clone() ExecuteCommand
	InitStopCatching()
	StopSignal()
	IsStopSignal() bool
	Execute(cp CommandProperties) error
}

// FullFileName function is used to construct full file name from parameters
func FullFileName(path string, fileName string) string {
	var fullName string
	if runtime.GOOS == "windows" {
		if fileName[0] == '\\' || strings.TrimPrefix(fileName, getRootPath(path)) != fileName {
			fullName = fileName
		} else {
			fullName = path + `\` + fileName
		}
	} else {
		if fileName[0] == '/' {
			fullName = fileName
		} else {
			fullName = path + "/" + fileName
		}
	}
	return fullName
}

// getRootPath function is used to extract root path from a valid path
func getRootPath(path string) string {
	if runtime.GOOS == "windows" {
		// path delimiter in Windows is \
		res := strings.SplitAfterN(path, `\`, 2)
		return res[0]
	}
	// path delimiter in Unix-like OS-es is /
	res := strings.SplitAfterN(path, `/`, 2)
	return res[0]
}

// checkRead function is very important - it reads from file, checking if there is a stop signal and also checking for error in reading
func checkRead(e ExecuteCommand, inputFile *os.File) (string, error) {
	if e.IsStopSignal() == true {
		return "", ErrStoppedExec
	}
	var buf = make([]byte, 1<<4)
	_, err := inputFile.Read(buf)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(buf), "\u0000"), nil
}

// checkWrite function is very important - it writes to file, checking if there is a stop signal and also checking for error in writing
func checkWrite(e ExecuteCommand, outputFile *os.File, text string) error {
	if e.IsStopSignal() == true {
		return ErrStoppedExec
	}
	n, err := outputFile.WriteString(text)
	if err != nil {
		return err
	}
	if n != len(text) {
		return fmt.Errorf("Wrote only %d characters of: %s", n, text)
	}
	return nil
}
