package main

import (
	"bufio"
	"fmt"
	"os"
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
		}
		fmt.Println(text)
	}
}
