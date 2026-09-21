package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		// command = command[:len(command)-2]
		command = strings.TrimSpace(command)
		builtins := []string{"exit", "echo", "type"}

		if command == "exit" {
			break
		} else if strings.HasPrefix(command, "echo ") {
			fmt.Println(command[5:])
		} else if strings.HasPrefix(command, "type ") {
			command = command[5:]
			if slices.Contains(builtins, command) {
				fmt.Println(command + " is a shell builtin")
			} else {
				fmt.Println(command + ": not found")
			}
		} else {
			fmt.Println(command, ": command not found")
		}
	}

}
