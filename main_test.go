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
