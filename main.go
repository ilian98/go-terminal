package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ilian98/go-terminal/commands"
	"github.com/ilian98/go-terminal/parser"
)

func main() {
	exitCommands := [...]string{"exit", "logout", "bye"}
	shellCommands := [...]string{"pwd", "cd"}

	path, err := os.Getwd()
	if err != nil {
		panic("Fatal error - cannot get current path!")
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
			for _, command := range shellCommands {
				if parsedCommand[0].Name == command {
					flagCommand = true
				}
			}
			if flagCommand == false {
				fmt.Println("No command with name: ", parsedCommand[0].Name)
				continue
			}

			command := commands.ExecuteCommand{
				Path:      path,
				Arguments: parsedCommand[0].Arguments,
				Options:   parsedCommand[0].Options,
				Input:     parsedCommand[0].Input,
				Output:    parsedCommand[0].Output,
			}
			commandName := strings.Title(parsedCommand[0].Name)
			runCommand := func() {
				result := reflect.ValueOf(&command).MethodByName(commandName).Call([]reflect.Value{})
				if r := result[0].Interface(); r != nil {
					fmt.Printf("%v\n", r.(error))
				}
			}
			if parsedCommand[0].BgRun == true {
				go runCommand()
			} else {
				runCommand()
				path = command.Path // path changed only when command is not run in bg mode
			}
		}
	}
}
