package commands

import (
	"fmt"
	"os"
	"testing"
)

func testingMkdir(t *testing.T, arguments []string, expectedErr string) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get working directory! - %w", err)
	}
	mkdir := Mkdir{}
	errMkdir := mkdir.Execute(CommandProperties{path, arguments, []string{}, os.Stdin, os.Stdout})
	defer func() {
		for _, argument := range arguments {
			os.RemoveAll(argument)
		}
	}()
	if errMkdir == nil {
		if expectedErr != "" {
			t.Errorf("Expected error %s, but got no error", expectedErr)
		}
		return
	}
	if errMkdir.Error() != expectedErr {
		t.Errorf("Expected error %s, but got: %v", expectedErr, errMkdir)
	}
}
func TestMkdir(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get working directory! - %w", err)
	}

	var tests = []struct {
		arguments []string
		err       string
	}{
		{[]string{}, ErrMkdirNoArgs.Error()},
		{[]string{"dir"}, ""},
		{[]string{"dir", "dir"}, FullFileName(path, "dir") + " - " + ErrMkdirExists.Error()},
		{[]string{"dir", "dir1"}, ""},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Mkdir test with arguments %v", test.arguments), func(t *testing.T) {
			testingMkdir(t, test.arguments, test.err)
		})
	}
}

func ExampleMkdir_Execute() {
	path, _ := os.Getwd()
	mkdir := Mkdir{}
	mkdir.Execute(CommandProperties{path, []string{"example-mkdir"}, []string{}, os.Stdin, os.Stdout})

	if err := os.RemoveAll(path + string(os.PathSeparator) + "example-mkdir"); err == nil {
		fmt.Println("example-mkdir directory was created!")
	}
	// Output:
	// example-mkdir directory was created!
}
