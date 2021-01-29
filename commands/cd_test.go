package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func testingCd(t *testing.T, cp CommandProperties, expectedResult string, expectedErr error) {
	cd := Cd{}
	err := cd.Execute(cp)
	if err != nil {
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected %v, but got: %w\n", expectedErr, err)
		}
		return
	}
	if expectedErr != nil {
		t.Errorf("Expected %v, but got no error\n", expectedErr)
		return
	}
	if resultPath := cd.GetPath(); resultPath != expectedResult {
		t.Errorf("Expected %s, but got: %s\n", expectedResult, resultPath)
	}
}
func TestCd(t *testing.T) {
	testPath, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get current path! - %w", err)
	}
	parentPath, err2 := filepath.Abs(testPath + `\..`)
	if err2 != nil {
		t.Fatal("Fatal error - cannot get parent path! - %w", err)
	}

	var tests = []struct {
		cp     CommandProperties
		result string
		err    error
	}{
		{CommandProperties{testPath, []string{"."}, []string{}, os.Stdin, os.Stdout}, testPath, nil},
		{CommandProperties{testPath, []string{".."}, []string{}, os.Stdin, os.Stdout}, parentPath, nil},
		{CommandProperties{testPath, []string{"..", "."}, []string{}, os.Stdin, os.Stdout}, "", ErrCdTooManyArgs},
		{CommandProperties{testPath, []string{"/not/existing/path"}, []string{}, os.Stdin, os.Stdout}, "", ErrCdPathNotExist},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Cd test with arguments %v ", test.cp.Arguments), func(t *testing.T) {
			testingCd(t, test.cp, test.result, test.err)
		})
	}
}

func ExampleCd_Execute() {
	cd := Cd{}
	cd.Execute(CommandProperties{"", []string{`\`}, []string{}, os.Stdin, os.Stdout})
	if cd.GetPath() == cd.getRootPath() {
		fmt.Printf("Terminal at root path!")
	}
	// Output:
	// Terminal at root path!
}
