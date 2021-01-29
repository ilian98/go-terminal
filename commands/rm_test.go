package commands

import (
	"fmt"
	"os"
	"testing"
)

func testingRm(t *testing.T, fileName string, dirName string, arguments []string, options []string, expectedErr string) {
	if dirName == "" {
		file, err := os.Create(fileName)
		if err != nil {
			t.Fatal("Fatal error - cannot make file! - %w", err)
		}
		file.Close()
		defer func() {
			os.Remove(fileName)
		}()
	} else {
		err := os.Mkdir(dirName, 0666)
		if err != nil {
			t.Fatal("Fatal error - cannot make directory! - %w", err)
		}
		file, err := os.Create(dirName + string(os.PathSeparator) + fileName)
		if err != nil {
			t.Fatal("Fatal error - cannot make file! - %w", err)
		}
		file.Close()
		defer func() {
			os.RemoveAll(dirName)
		}()
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get working directory! - %w", err)
	}
	rm := Rm{}
	errRm := rm.Execute(CommandProperties{wd, arguments, options, os.Stdin, os.Stdout})
	if errRm == nil {
		if expectedErr != "" {
			t.Errorf("Expected error %s, but got no error", expectedErr)
		}
		return
	}
	if errRm.Error() != expectedErr {
		t.Errorf("Expected error %s, but got: %v", expectedErr, errRm)
	}
}
func TestRm(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get working directory! - %w", err)
	}

	var tests = []struct {
		fileName  string
		dirName   string
		arguments []string
		options   []string
		err       string
	}{
		{"file", "", []string{}, []string{}, ErrRmNoArgs.Error()},
		{"f", "dir", []string{"dir"}, []string{}, FullFileName(path, "dir") + " " + ErrRmIsDir.Error()},
		{"file", "", []string{"file"}, []string{}, ""},
		{"file", "", []string{"file1"}, []string{}, FullFileName(path, "file1") + " " + ErrRmInvalidName.Error()},
		{"file", "", []string{"file", "file"}, []string{}, FullFileName(path, "file") + " " + ErrRmInvalidName.Error()},
		{"f", "dir", []string{"dir"}, []string{"r"}, ""},
		{"f", "dir", []string{"dir/f"}, []string{"r"}, FullFileName(path, "dir/f") + " " + ErrRmIsFile.Error()},
		{"f", "dir", []string{"dir", "dir/f"}, []string{"r"}, FullFileName(path, "dir/f") + " " + ErrRmInvalidName.Error()},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Rm test with fileName %s in dirName %s, arguments %v and options %v", test.fileName, test.dirName, test.arguments, test.options), func(t *testing.T) {
			testingRm(t, test.fileName, test.dirName, test.arguments, test.options, test.err)
		})
	}
}

func ExampleRm_Execute() {
	path, _ := os.Getwd()
	file, _ := os.Create(path + string(os.PathSeparator) + "not-existing-file")
	file.Close()
	rm := Rm{}
	rm.Execute(CommandProperties{path, []string{"not-existing-file"}, []string{}, os.Stdin, os.Stdout})

	if err := os.Remove(path + string(os.PathSeparator) + "not-existing-file"); err != nil {
		fmt.Println("not-existing-file removed!")
	}
	// Output:
	// not-existing-file removed!
}
