package commands

import (
	"os"
	"testing"
)

func TestOpenInputOutputFile(t *testing.T) {
	cp1 := CommandProperties{"proba", []string{}, []string{}, "non-existing-file", ""}
	file1, _, err1 := cp1.openInputOutputFiles()
	if err1 == nil {
		t.Error("Expecting error\n")
	}
	if file1 != nil {
		t.Error("Expecting no file\n")
	}

	cp2 := CommandProperties{"proba", []string{}, []string{}, "", ""}
	file2, file3, err2 := cp2.openInputOutputFiles()
	if err2 != nil {
		t.Error("Expecting no error\n")
	}
	if file2 != os.Stdin {
		t.Error("Expecting file to be stdin\n")
	}
	if file3 != os.Stdout {
		t.Error("Expecting file to be stdout\n")
	}
}
