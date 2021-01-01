package parser

import (
	"fmt"
	"testing"
)

func TestExitParse(t *testing.T) {
	parsedCommand := Parse("exit\n")
	if parsedCommand[0].Name != "exit" {
		t.Error("Name of exit command not parsed correctly!")
	}
}

func ExampleParse() {
	parsedCommand := Parse("exit\n")
	fmt.Println(parsedCommand[0].Name)
	// Output:
	// exit
}
