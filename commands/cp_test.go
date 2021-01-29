package commands

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func testingCp(t *testing.T, source string, arguments []string, expectedErr error) {
	if source != "" {
		file, err := os.Create(source)
		if err != nil {
			t.Fatal("Fatal error - cannot make file! - %w", err)
		}
		file.Close()
		defer func() {
			os.Remove(source)
		}()
	}

	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get working directory! - %w", err)
	}
	cp := Cp{}
	errCp := cp.Execute(CommandProperties{path, arguments, []string{}, os.Stdin, os.Stdout})
	if errCp == nil {
		if expectedErr != nil {
			t.Errorf("Expected error %v, but got no error", expectedErr)
		} else {
			os.Remove(arguments[1])
		}
		return
	}
	if !errors.Is(errCp, expectedErr) {
		t.Errorf("Expected error %v, but got: %v", expectedErr, errCp)
	}
}
func TestCp(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get working directory! - %w", err)
	}

	var tests = []struct {
		source    string
		arguments []string
		err       error
	}{
		{"", []string{"test1"}, ErrCpTwoArgs},
		{"", []string{"not-existing-file", "test2"}, ErrCpInvalidName},
		{"", []string{"test3", "test3"}, ErrCpSame},
		{"test4", []string{"test4", "test4'"}, nil},
		{"", []string{getRootPath(path), "test5"}, ErrCpIsDir},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Rm test with arguments %v", test.arguments), func(t *testing.T) {
			testingCp(t, test.source, test.arguments, test.err)
		})
	}
}

func ExampleCp_Execute() {
	path, _ := os.Getwd()
	file, _ := os.Create(path + string(os.PathSeparator) + "not-existing-file")
	file.Close()
	defer os.Remove(path + string(os.PathSeparator) + "not-existing-file")
	cp := Cp{}
	cp.Execute(CommandProperties{path, []string{"not-existing-file", "copy"}, []string{}, os.Stdin, os.Stdout})

	if err := os.Remove(path + string(os.PathSeparator) + "copy"); err == nil {
		fmt.Println("not-existing-file was copied!")
	}
	// Output:
	// not-existing-file was copied!
}
