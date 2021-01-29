package commands

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// CommandProperties is used for storing the properties of command that will be executed
type CommandProperties struct {
	Path       string //  Path is used for storing the current path in the terminal for this command
	Arguments  []string
	Options    []string
	InputFile  *os.File // InputFile is used for reading the input it could be stdin
	OutputFile *os.File // OutputFile is used for reading the input it could be stdin
}

// ExecuteCommand is interface for executing commands
type ExecuteCommand interface {
	GetName() string
	GetPath() string
	Execute(cp CommandProperties) error
	Clone() ExecuteCommand
}

// CloseInputOutputFiles is used by Exectute to close the opened input and output files
func CloseInputOutputFiles(inputFile *os.File, outputFile *os.File) {
	if inputFile != os.Stdin {
		inputFile.Close()
	}
	if outputFile != os.Stdout {
		outputFile.Close()
	}
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

func checkWrite(outputFile *os.File, text string) error {
	n, err := outputFile.WriteString(text)
	if err != nil {
		return err
	}
	if n != len(text) {
		return fmt.Errorf("Wrote only %d characters of: %s", n, text)
	}
	return nil
}
