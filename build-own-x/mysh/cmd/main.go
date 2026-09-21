package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		command = command[:len(command)-2]

		if command == "exit" {
			break
		}
		fmt.Println(command, ": command not found")
	}

}
