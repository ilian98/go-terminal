package commands

import (
	"fmt"
	"os"
	"testing"
)

func testingCat(t *testing.T, input string, inputText string, arguments []string, expectedResult string, expectedErr string) {
	inputR, inputW, err := os.Pipe()
	if err != nil {
		t.Fatal("Fatal error - cannot make pipe! - %w", err)
	}
	stdin := os.Stdin
	defer func() { os.Stdin = stdin }()
	os.Stdin = inputR

	outputR, outputW, err := os.Pipe()
	if err != nil {
		t.Fatal("Fatal error - cannot make pipe! - %w", err)
	}
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	os.Stdout = outputW

	inputW.WriteString(inputText)
	inputW.Close()
	cat := Cat{}
	if err := cat.Execute(CommandProperties{"test/path", arguments, []string{}, input, ""}); err != nil {
		if expectedErr == "" {
			t.Errorf("Expected no error, but got: %v", err)
		} else if err.Error() != expectedErr {
			t.Errorf("Exepected error %s, but got: %v", expectedErr, err)
		}
		return
	}
	if expectedErr != "" {
		t.Errorf("Expected error %s, but got no error", expectedErr)
		return
	}

	output := make([]byte, len(expectedResult))
	if _, err := outputR.Read(output); err != nil {
		t.Fatal("Fatal error - cannot read from pipe! - %w", err)
	}
	if string(output) != expectedResult {
		t.Errorf("Expecting %s, but got: %s", expectedResult, string(output))
	}
}
func TestCat(t *testing.T) {
	var tests = []struct {
		input     string
		inputText string
		arguments []string
		result    string
		err       string
	}{
		{"", "first test", []string{}, "first test", ""},
		{"", "second test", []string{"no-file"}, "", "no-file - file does not exist"},
		{"", "third test", []string{"no-file1", "no-file2"}, "", "no-file1 - file does not exist\nno-file2 - file does not exist"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Cat test with input %s, inputText %s and arguments %v ", test.input, test.inputText, test.arguments), func(t *testing.T) {
			testingCat(t, test.input, test.inputText, test.arguments, test.result, test.err)
		})
	}
}

func ExampleCat_Execute() {
	path, _ := os.Getwd()
	file, _ := os.OpenFile(path+string(os.PathSeparator)+"example-file", os.O_CREATE|os.O_WRONLY, 0666)
	defer func() {
		os.Remove(path + string(os.PathSeparator) + "example-file")
	}()
	file.WriteString("cat command example")
	file.Close()

	cat := Cat{}
	cat.Execute(CommandProperties{path, []string{"example-file"}, []string{}, "", ""})
	// Output:
	// cat command example
}