package commands

import (
	"os"
	"strings"
	"testing"
)

func TestLs(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get current path! - %w", err)
	}
	path += string(os.PathSeparator) + "example-dir"
	if err := os.Mkdir(path, 0666); err != nil {
		t.Fatal("Fatal error - cannot make directory in current path! - %w", err)
	}
	pathFile := path + string(os.PathSeparator) + "new-file"
	file, err := os.OpenFile(pathFile, os.O_CREATE, 0666)
	if err != nil {
		t.Fatal("Fatal error - cannot make file in new directory new-file! - %w", err)
	}
	defer func() {
		file.Close()
		os.Remove(pathFile)
		os.Remove(path)
	}()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Fatal error - cannot make pipe! - %w", err)
	}

	ls := Ls{}
	if err := ls.Execute(CommandProperties{path, []string{}, []string{"l"}, os.Stdin, w}); err != nil {
		t.Errorf("Expecting no error from Ls function, but got: %w\n", err)
		return
	}

	output := make([]byte, 1<<10)
	if _, err := r.Read(output); err != nil {
		t.Fatal("Fatal error - cannot read from pipe! - %w", err)
	}
	parts := strings.Split(strings.SplitN(string(output), "\n", 2)[0], " ")
	if len(parts) != 6 { // format of ls -l should have 6 columns
		t.Errorf("Expecting 6 columns, but got: %s", output)
		return
	}
	if parts[0] != "-rw-rw-rw-" { // default open mode for file is -rw-rw-rw-
		t.Errorf("Expecting -rw-rw-rw, but got: %s", parts[0])
		return
	}
	if parts[1] != "0" {
		t.Errorf("Expecting empty file, but file is with size %s", parts[1])
		return
	}
	// we only check the length of parts[2:4] which should be respectively in format hh:mm dd mmm
	if len(parts[2]) != 5 || len(parts[3]) != 2 || len(parts[4]) != 3 {
		t.Errorf("Expecting hh:mm dd mmm format of time and data, but got: %s %s %s", parts[2], parts[3], parts[4])
		return
	}
	if parts[5] != "new-file" {
		t.Errorf("Expecting file name to be new-file, but got: %s", parts[5])
	}
}

func ExampleLs_Execute() {
	path, _ := os.Getwd()
	path += string(os.PathSeparator) + "example-dir"
	os.Mkdir(path, 0666)
	pathFile := path + string(os.PathSeparator) + "new-file"
	file, _ := os.OpenFile(pathFile, os.O_CREATE, 0666)
	defer func() {
		file.Close()
		os.Remove(pathFile)
		os.Remove(path)
	}()

	ls := Ls{}
	ls.Execute(newCp(path, []string{}, []string{}))
	// Output:
	// new-file
}
