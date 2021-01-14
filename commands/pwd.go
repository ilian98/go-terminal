package commands

import (
	"fmt"
)

// Pwd function is go implementation of pwd command
func Pwd(arguments []string, options []string, input string, output string) (func() error, error) {
	inputFile, outputFile, err := openInputOutputFiles(input, output)
	if err != nil {
		return nil, err
	}

	return func() error {
		defer closeInputOutputFiles(input, inputFile, output, outputFile)
		n, err := outputFile.WriteString(Path)
		if err != nil {
			return err
		}
		if n != len(Path) {
			return fmt.Errorf("Wrote only %d characters of current path: %s", n, Path)
		}

		return nil
	}, nil
}
