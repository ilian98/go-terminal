package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ilian98/go-terminal/commands"
	"github.com/ilian98/go-terminal/parser"
)

func main() {
	exitCommands := [...]string{"exit", "logout", "bye"}

	path, err := os.Getwd()
	if err != nil {
		panic("Fatal error - cannot get current path!")
	}
	commands.Path = path

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("")
		fmt.Println(commands.Path)
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

			var commandFunc func() error
			if parsedCommand[0].Name == "pwd" {
				f, err := commands.Pwd(parsedCommand[0].Arguments, parsedCommand[0].Options, parsedCommand[0].Input, parsedCommand[0].Output)
				if err != nil {
					fmt.Printf("%v\n", err)
					continue
				}
				commandFunc = f
			} else {
				fmt.Println("No command with name: ", parsedCommand[0].Name)
				continue
			}

			if parsedCommand[0].BgRun == true {
				go func() {
					err = commandFunc()
					if err != nil {
						fmt.Printf("%v\n", err)
					}
				}()
			} else {
				err = commandFunc()
				if err != nil {
					fmt.Printf("%v\n", err)
					continue
				}
			}

		}
	}
}
