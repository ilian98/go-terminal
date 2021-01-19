package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// ErrTooManyArgs indicates that the cd command has more than 1 argument
	ErrTooManyArgs = errors.New("Too many arguments")
	// ErrPathLeadsToFile indicates that the path in cd command is to file
	ErrPathLeadsToFile = errors.New("path leads to file, not directory")
	// ErrPathNotExist indicates that the path in cd command does not exist
	ErrPathNotExist = errors.New("path does not exist")
)

// Cd is method of structure ExecuteCommand that executes go implementation of cd command
func (e *ExecuteCommand) Cd() error {
	inputFile, outputFile, err := e.openInputOutputFiles(e.Input, e.Output)
	defer e.closeInputOutputFiles(e.Input, inputFile, e.Output, outputFile)

	if err != nil {
		return err
	}

	if len(e.Arguments) == 0 {
		e.Path = getRootPath(e.Path)
		return nil
	}
	if len(e.Arguments) > 1 {
		return ErrTooManyArgs
	}

	path := e.Arguments[0]
	if len(path) == 0 {
		e.Path = getRootPath(e.Path)
		return nil
	}

	var tryPath string
	if runtime.GOOS == "windows" {
		if path[0] == '\\' || strings.TrimPrefix(path, getRootPath(e.Path)) != path {
			tryPath = path
		} else {
			tryPath = e.Path + `\` + path
		}
	} else {
		if path[0] == '/' {
			tryPath = path
		} else {
			tryPath = e.Path + "/" + path
		}
	}

	stat, err := os.Stat(tryPath)
	if err == nil && stat.IsDir() {
		p, err := filepath.Abs(tryPath)
		if err != nil {
			return err
		}
		e.Path = p
		return nil
	} else if err == nil {
		return fmt.Errorf("%s - %w", tryPath, ErrPathLeadsToFile)
	} else if os.IsNotExist(err) {
		return fmt.Errorf("%s - %w", tryPath, ErrPathNotExist)
	}
	return err
}

func getRootPath(path string) string {
	if runtime.GOOS == "windows" {
		// path delimiter in Windows is \
		res := strings.SplitAfterN(path, `\`, 2)
		if len(res) == 0 {
			panic("Cannot get root path!")
		}
		return res[0]
	}
	// path delimiter in Unix-like OS-es is /
	res := strings.SplitAfterN(path, `/`, 2)
	if len(res) == 0 {
		panic("Cannot get root path!")
	}
	return res[0]
}
