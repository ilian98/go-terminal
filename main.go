package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ilian98/go-terminal/parser"
)

func main() {
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
		if parsedCommand[0].Name == "exit" {
			break
		}
		fmt.Println("")
	}
}
