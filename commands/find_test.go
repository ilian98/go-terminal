package commands

import (
	"os"
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get current path! - %w", err)
	}
	pathFile := path + string(os.PathSeparator) + "new-file"
	file, err := os.Create(pathFile)
	if err != nil {
		t.Fatal("Fatal error - cannot make file in new directory new-file! - %w", err)
	}
	defer func() {
		file.Close()
		os.Remove(pathFile)
	}()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Fatal error - cannot make pipe! - %w", err)
	}

	find := Find{}
	if err := find.Execute(CommandProperties{path, []string{"new-file", "new-file"}, []string{}, os.Stdin, w}); err != nil {
		t.Errorf("Expecting no error from Find function, but got: %w\n", err)
		return
	}

	output := make([]byte, 1<<10)
	if _, err := r.Read(output); err != nil {
		t.Fatal("Fatal error - cannot read from pipe! - %w", err)
	}
	s := strings.Split(string(output), "\n")
	if len(s) < 2 {
		t.Errorf("Expecting 2 lines, but got: %s", output)
		return
	}
	expectedResult := "new-file found - " + pathFile
	if s[0] != expectedResult || s[1] != expectedResult {
		t.Errorf("Expecting %s, but got: %s and %s", expectedResult, s[0], s[1])
		return
	}

	if err := find.Execute(CommandProperties{path, []string{}, []string{}, os.Stdin, w}); err != ErrFindNoArgs {
		t.Errorf("Expecting error %v, but got: %w\n", ErrFindNoArgs, err)
		return
	}
}

func ExampleFind_Execute() {
	path, _ := os.Getwd()
	find := Find{}
	find.Execute(newCp(path, []string{"not-existing-file"}, []string{}))
	// Output:
	// not-existing-file not found
}
