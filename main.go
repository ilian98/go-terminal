package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ilian98/go-terminal/parser"
)

func main() {
	exitCommands := [...]string{"exit", "logout", "bye"}

	path, err := os.Getwd()
	if err != nil {
		panic("Fatal error - cannot get current path!")
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(path)
		fmt.Print("$ ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Couldn't read command!")
			continue
		}
		parsedCommand := parser.Parse(text)
		for _, exitCommand := range exitCommands {
			if parsedCommand[0].Name == exitCommand {
				if exitCommand == "bye" {
					fmt.Println("bye")
				}
				return
			}
		}
		fmt.Println("")
	}
}
