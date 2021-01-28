package commands

import (
	"os"
	"testing"
)

func TestPwd(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Fatal error - cannot make pipe! - %w", err)
	}

	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = w

	testPath := "testPwd"
	pwd := Pwd{}
	if err := pwd.Execute(CommandProperties{testPath, []string{}, []string{}, "", ""}); err != nil {
		t.Error("Expecting no error from Pwd function\n")
	}

	output := make([]byte, len(testPath))

	if _, err := r.Read(output); err != nil {
		t.Fatal("Fatal error - cannot read from pipe! - %w", err)
	}
	if string(output) != testPath {
		t.Errorf("Expecting %s, but got: %s", testPath, string(output))
	}
}

func ExamplePwd_Execute() {
	pwd := Pwd{}
	pwd.Execute(CommandProperties{"Example/Path", []string{}, []string{}, "", ""})
	// Output:
	// Example/Path
}
