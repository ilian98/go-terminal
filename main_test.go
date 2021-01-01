package main

import (
	"os"
	"testing"
	"time"
)

func TestExitCommand(t *testing.T) {
	input := []byte("exit\n")
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

	mainFinished := false
	go func() {
		main()
		mainFinished = true
	}()

	time.Sleep(100 * time.Millisecond)
	if mainFinished == false {
		t.Error("Exit command not working!")
	}
}
