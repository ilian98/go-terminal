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
		I.RegisterExitCommand(exitCommand)
	}

	commands := [...]commands.ExecuteCommand{
		&commands.Pwd{}, &commands.Cd{}, &commands.Ls{}, &commands.Cat{}, &commands.Ping{},
	}
	for _, command := range commands {
		I.RegisterCommand(command)
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
		if len(parsedCommand) == 1 {
			status := I.ExecuteCommand(parsedCommand[0])
			if status == interpreter.ExitCommand {
				if parsedCommand[0].Name == "bye" {
					fmt.Println("bye")
				}
				fmt.Println("")
				return
			}
			if status == interpreter.InvalidCommandName {
				fmt.Printf("No command with name: %s\n", parsedCommand[0].Name)
			}
		}
	}
}
