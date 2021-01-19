package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func testingCd(t *testing.T, command ExecuteCommand, expectedResult string, expectedErr error) {
	err := command.Cd()
	if err != nil {
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected %v, but got: %v\n", expectedErr, errors.Unwrap(err))
		}
		return
	}
	if expectedErr != nil {
		t.Errorf("Expected %v, but got no error\n", expectedErr)
	}
	if command.Path != expectedResult {
		t.Errorf("Expected %s, but got: %s\n", expectedResult, command.Path)
	}
}
func TestCd(t *testing.T) {
	testPath, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get current path!")
	}
	parentPath, err2 := filepath.Abs(testPath + `\..`)
	if err2 != nil {
		t.Fatal("Fatal error - cannot get parent path!")
	}

	var tests = []struct {
		command ExecuteCommand
		result  string
		err     error
	}{
		{ExecuteCommand{testPath, []string{"."}, []string{}, "", ""}, testPath, nil},
		{ExecuteCommand{testPath, []string{".."}, []string{}, "", ""}, parentPath, nil},
		{ExecuteCommand{testPath, []string{"..", "."}, []string{}, "", ""}, "", ErrTooManyArgs},
		{ExecuteCommand{testPath, []string{"/not/existing/path"}, []string{}, "", ""}, "", ErrPathNotExist},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Cd test with arguments %v ", test.command.Arguments), func(t *testing.T) {
			testingCd(t, test.command, test.result, test.err)
		})
	}
}

func ExampleExecuteCommand_Cd() {
	e := ExecuteCommand{"", []string{`\`}, []string{}, "", ""}
	e.Cd()
	if e.Path == getRootPath(e.Path) {
		fmt.Printf("Terminal at root path!")
	}
	// Output:
	// Terminal at root path!
}
