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
	// ErrCdTooManyArgs indicates that the cd command has more than 1 argument
	ErrCdTooManyArgs = errors.New("Too many arguments")
	// ErrCdPathLeadsToFile indicates that the path in cd command is to file
	ErrCdPathLeadsToFile = errors.New("path leads to file, not directory")
	// ErrCdPathNotExist indicates that the path in cd command does not exist
	ErrCdPathNotExist = errors.New("path does not exist")
)

// Cd is a structure for cd command, implementing ExecuteCommand interface
type Cd struct {
	path string
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

// Execute is go implementation of cd command
func (c *Cd) Execute(cp CommandProperties) error {
	c.path = cp.Path
	inputFile, outputFile := cp.InputFile, cp.OutputFile
	defer CloseInputOutputFiles(inputFile, outputFile)

	if len(cp.Arguments) == 0 {
		c.path = c.getRootPath()
		return nil
	}
	if len(cp.Arguments) > 1 {
		return ErrCdTooManyArgs
	}

	path := cp.Arguments[0]
	if len(path) == 0 {
		c.path = c.getRootPath()
		return nil
	}

	var tryPath string
	if runtime.GOOS == "windows" {
		if path[0] == '\\' || strings.TrimPrefix(path, c.getRootPath()) != path {
			tryPath = path
		} else {
			tryPath = c.path + `\` + path
		}
	} else {
		if path[0] == '/' {
			tryPath = path
		} else {
			tryPath = c.path + "/" + path
		}
	}

	stat, err := os.Stat(tryPath)
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

func (c *Cd) getRootPath() string {
	if runtime.GOOS == "windows" {
		// path delimiter in Windows is \
		res := strings.SplitAfterN(c.path, `\`, 2)
		if len(res) == 0 {
			panic("Cannot get root path!")
		}
		return res[0]
	}
	// path delimiter in Unix-like OS-es is /
	res := strings.SplitAfterN(c.path, `/`, 2)
	if len(res) == 0 {
		panic("Cannot get root path!")
	}
	return res[0]
}
