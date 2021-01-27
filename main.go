package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ilian98/go-terminal/commands"
	"github.com/ilian98/go-terminal/parser"
)

func runCommand(command commands.ExecuteCommand, cp commands.CommandProperties) string {
	if err := command.Execute(cp); err != nil {
		fmt.Printf("%v\n", err)
	}
	return command.GetPath()
}

func main() {
	exitCommands := [...]string{"exit", "logout", "bye"}
	shellCommandsName := [...]string{"pwd", "cd"}

	path, err := os.Getwd()
	if err != nil {
		panic("Fatal error - cannot get current path!")
	}
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("")
		fmt.Println(path)
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
			// check if command is for exiting the terminal
			for _, exitCommand := range exitCommands {
				if parsedCommand[0].Name == exitCommand {
					if exitCommand == "bye" {
						fmt.Println("bye")
					}
					fmt.Println("")
					return
				}
			}

			flagCommand := false
			var indCommand int
			for i, command := range shellCommandsName {
				if parsedCommand[0].Name == command {
					flagCommand = true
					indCommand = i
				}
			}
			if flagCommand == false {
				fmt.Printf("No command with name: %s\n", parsedCommand[0].Name)
				continue
			}

			cp := commands.CommandProperties{
				Path:      path,
				Arguments: parsedCommand[0].Arguments,
				Options:   parsedCommand[0].Options,
				Input:     parsedCommand[0].Input,
				Output:    parsedCommand[0].Output,
			}

			shellCommandsExecute := [...]commands.ExecuteCommand{&commands.Pwd{}, &commands.Cd{}}
			if parsedCommand[0].BgRun == true {
				go runCommand(shellCommandsExecute[indCommand], cp)
			} else {
				path = runCommand(shellCommandsExecute[indCommand], cp) // path changed only when command is not run in bg mode
			}
		}
	}
}
