package commands

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func testingMv(t *testing.T, source string, arguments []string, expectedErr error) {
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
	mv := Mv{}
	errMv := mv.Execute(newCp(path, arguments, []string{}))
	if errMv == nil {
		if expectedErr != nil {
			t.Errorf("Expected error %v, but got no error", expectedErr)
		} else {
			os.Remove(arguments[1])
		}
		return
	}
	if !errors.Is(errMv, expectedErr) {
		t.Errorf("Expected error %v, but got: %v", expectedErr, errMv)
	}
}
func TestMv(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get working directory! - %w", err)
	}

	var tests = []struct {
		source    string
		arguments []string
		err       error
	}{
		{"", []string{"test1"}, ErrMvTwoArgs},
		{"", []string{"not-existing-file", "test2"}, ErrMvInvalidName},
		{"", []string{"test3", "test3"}, ErrMvSame},
		{"test4", []string{"test4", "test4'"}, nil},
		{"", []string{getRootPath(path), "test5"}, ErrMvIsDir},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Rm test with arguments %v", test.arguments), func(t *testing.T) {
			testingMv(t, test.source, test.arguments, test.err)
		})
	}
}

func ExampleMv_Execute() {
	path, _ := os.Getwd()
	file, _ := os.Create(path + string(os.PathSeparator) + "not-existing-file")
	file.Close()
	mv := Mv{}
	mv.Execute(newCp(path, []string{"not-existing-file", "example-mv"}, []string{}))

	if err := os.Remove(path + string(os.PathSeparator) + "example-mv"); err == nil {
		fmt.Println("not-existing-file was renamed to example-mv!")
	}
	// Output:
	// not-existing-file was renamed to example-mv!
}
