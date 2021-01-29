package interpreter

import (
	"testing"

	"github.com/ilian98/go-terminal/commands"
	"github.com/ilian98/go-terminal/parser"
)

func TestRegisterExitCommand(t *testing.T) {
	var i Interpreter
	if err := i.RegisterExitCommand("exit"); err != nil {
		t.Error("Expecting no error\n")
	}
	if len(i.exitCommands) != 1 || i.exitCommands[0] != "exit" {
		t.Errorf("Expecting %v, but got %v\n", []string{"exit"}, i.exitCommands)
	}

	if err := i.RegisterExitCommand("exit"); err != ErrCommandExists {
		t.Error("Expecting error %w, but got: \n", ErrCommandExists, err)
	}

	if err := i.RegisterExitCommand("exit2"); err != nil {
		t.Error("Expecting no error\n")
	}
	if len(i.exitCommands) != 2 || i.exitCommands[1] != "exit2" {
		t.Errorf("Expecting %v, but got %v\n", []string{"exit", "exit2"}, i.exitCommands)
	}
}

func TestRegisterCommand(t *testing.T) {
	var i Interpreter
	if err := i.RegisterCommand(&commands.Pwd{}); err != nil {
		t.Error("Expecting no error\n")
	}
	if len(i.shellCommandsName) != 1 || i.shellCommandsName[0] != "pwd" {
		t.Errorf("Expecting %v, but got %v\n", []string{"pwd"}, i.shellCommandsName)
	}

	if err := i.RegisterCommand(&commands.Pwd{}); err != ErrCommandExists {
		t.Error("Expecting error %w, but got: \n", ErrCommandExists, err)
	}

	if err := i.RegisterCommand(&commands.Cd{}); err != nil {
		t.Error("Expecting no error\n")
	}
	if len(i.shellCommandsName) != 2 || i.shellCommandsName[1] != "cd" {
		t.Errorf("Expecting %v, but got %v\n", []string{"pwd", "cd"}, i.shellCommandsName)
	}
}

func ExampleInterpreter() {
	var i Interpreter
	i.RegisterCommand(&commands.Pwd{})
	i.Path = "example/path"
	i.InterpretCommand([]parser.Command{
		{
			Name:      "pwd",
			Arguments: []string{},
			Options:   []string{},
			Input:     "",
			Output:    "",
			BgRun:     false,
		},
	})

	// Output:
	// example/path
}
