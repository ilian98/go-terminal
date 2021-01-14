package commands

import (
	"os"
	"testing"
)

func TestOpenInputOutputFile(t *testing.T) {
	file1, _, err1 := openInputOutputFiles("non-existing-file", "")
	if err1 == nil {
		t.Error("Expecting error\n")
	}
	if file1 != nil {
		t.Error("Expecting no file\n")
	}

	file2, file3, err2 := openInputOutputFiles("", "")
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
