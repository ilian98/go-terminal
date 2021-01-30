package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func testingExitCommand(t *testing.T, exitCommand string) {
	input := []byte(exitCommand + "\n")
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write(input)
	if err != nil {
		t.Error(err)
	}

	stdin := os.Stdin
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	finishChannel := make(chan struct{}, 1)
	go func() {
		main()
		finishChannel <- struct{}{}
	}()

	select {
	case <-finishChannel:
	case <-time.After(100 * time.Millisecond):
		t.Error("Exit command not working!")
	}
}
func TestExitCommand(t *testing.T) {
	exitCommands := [...]string{"exit", "logout", "bye"}

	for _, exitCommand := range exitCommands {
		t.Run(fmt.Sprintf("Еxit command %s", exitCommand), func(t *testing.T) {
			testingExitCommand(t, exitCommand)
		})
	}
}

func testingPipe(t *testing.T, inputText string, fileData string, expectedResult string) {
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

	inputText += "exit\n"
	if _, err := inputW.WriteString(inputText); err != nil {
		t.Fatal("Fatal error - cannot write to pipe! - %w", err)
	}
	inputW.Close()

	file, err := os.OpenFile("TestPipe", os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal("Fatal error - cannot open file! - %w", err)
	}
	file.WriteString(fileData)
	file.Close()

	finish := make(chan struct{}, 1)
	go func() {
		main()
		finish <- struct{}{}
	}()
	select {
	case <-finish:
	case <-time.After(6 * time.Second):
		t.Error("Main didn't finish for 6 seconds")
		<-finish
		return
	}
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get working directory! - %w", err)
	}
	expectedResult = "\n" + path + "\n$ " + expectedResult + "\n" + path + "\n$ "
	output := make([]byte, len(expectedResult))
	if _, err := outputR.Read(output); err != nil {
		t.Fatal("Fatal error - cannot read from pipe! - %w", err)
	}
	if string(output) != expectedResult {
		t.Errorf("Expecting %s, but got: %s.", expectedResult, string(output))
	}
}
func TestPipe(t *testing.T) {
	testPath, err := os.Getwd()
	if err != nil {
		t.Fatal("Fatal error - cannot get current path! - %w", err)
	}

	var tests = []struct {
		inputText string
		fileData  string
		result    string
	}{
		{"cat TestPipe\n", "test1", "test1"},
		{"pwd | cat > TestPipe | cat\n", "", ""},
		{"pwd | cat\n", "", testPath},
		{"cat TestPipe TestPipe | cat | cat | cat | cat\n", "test4", "test4test4"},
		{"ping noibg.com | cd\n", "", "write |1: The pipe is being closed.\n"},
		{"pwd | ls -l | cmd1 | cd\n", "", "No command with name: cmd1\n"},
	}

	file, err := os.Create("TestPipe")
	file.Close()
	if err != nil {
		t.Fatal("Fatal error - cannot make file for testing! - %w", err)
	}
	defer func() {
		os.Remove("TestPipe")
	}()
	for _, test := range tests {
		t.Run(fmt.Sprintf("Test main with inputText %s and fileData %s", test.inputText, test.fileData), func(t *testing.T) {
			testingPipe(t, test.inputText, test.fileData, test.result)
		})
	}
}
