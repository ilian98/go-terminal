package commands

import (
	"os"
	"testing"
)

func TestPwd(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = w

	Path = "testPwd"
	pwd, err := Pwd([]string{}, []string{}, "", "")
	if err != nil {
		t.Error("Expecting no error from Pwd function\n")
	}
	if err := pwd(); err != nil {
		t.Error("Expecting no error from pwd command\n")
	}

	output := make([]byte, len(Path))

	if _, err := r.Read(output); err != nil {
		t.Fatal(err)
	}
	if string(output) != Path {
		t.Errorf("Expecting %s, but got: %s", Path, string(output))
	}
}
