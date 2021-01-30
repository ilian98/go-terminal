package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ilian98/go-terminal/commands"
	"github.com/ilian98/go-terminal/interpreter"
	"github.com/ilian98/go-terminal/parser"
)

// I is the interpreter which is used by main to operate with the parsed commands
var I interpreter.Interpreter

func init() {
	for _, exitCommand := range []string{"exit", "logout", "bye"} {
		if err := I.RegisterExitCommand(exitCommand); err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	commands := [...]commands.ExecuteCommand{
		&commands.Pwd{}, &commands.Cd{}, &commands.Ls{}, &commands.Cat{},
		&commands.Cp{}, &commands.Mv{}, &commands.Mkdir{}, &commands.Rm{},
		&commands.Find{}, &commands.Ping{},
	}
	for _, command := range commands {
		if err := I.RegisterCommand(command); err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic("Fatal error - cannot get current path!")
	}
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
	I.Path = path

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("")
		fmt.Println(I.Path)
		fmt.Print("$ ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Couldn't read command!")
			continue
		}
		parsedCommand, err := parser.Parse(text)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		statuses := I.InterpretCommand(parsedCommand)
		if len(statuses) == 1 && statuses[0].Code == interpreter.ExitCommand {
			if parsedCommand[0].Name == "bye" {
				fmt.Println("bye")
			}
			fmt.Println("")
			return
		}
		for _, status := range statuses {
			if status.Code == interpreter.InvalidCommandName {
				fmt.Printf("No command with name: %s\n", status.Command)
			}
		}
	}
}
