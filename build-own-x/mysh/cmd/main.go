package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

		command = strings.TrimSpace(command)
		builtins := map[string]string{
			"exit":   "exit",
			"echo":   "echo",
			"type":   "type",
			"ls":     "Get-ChildItem",
			"pwd":    "Get-Location",
			"cd":     "Set-Location",
			"go":     "go",
			"git":    "git",
			"docker": "docker",
		}

		if command == "exit" {
			break
		} else if strings.HasPrefix(command, "echo ") {
			fmt.Println(command[5:])
		} else if strings.HasPrefix(command, "ls") {
			dirs, err := os.ReadDir(".")
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error reading directory:", err)
				os.Exit(1)
			}
			for _, dir := range dirs {
				fmt.Println(dir.Name())
			}

		} else if strings.HasPrefix(command, "pwd") {
			location, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting current directory:", err)
				os.Exit(1)
			}
			fmt.Println(location)

		} else if strings.HasPrefix(command, "git") {
			args := strings.Fields(command)
			cmd := exec.Command(builtins["git"], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "git error:", err)
			}
		} else if strings.HasPrefix(command, "docker") {
			args := strings.Fields(command)
			cmd := exec.Command(builtins["docker"], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "docker error:", err)
			}
		} else if strings.HasPrefix(command, "cd ") {
			dir := command[3:]
			err := os.Chdir(dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error changing directory:", err)
				os.Exit(1)
			}

		} else if strings.HasPrefix(command, "type ") {
			command = command[5:]
			if builtins[command] != "" {

				path, err := exec.LookPath(builtins[command])

				if err != nil {
					fmt.Fprintln(os.Stderr, "Error looking up command:", err)
					os.Exit(1)
				}
				fmt.Println(command+" is ", path)
				// fmt.Printf("fortune is available at %s\n", path)
			} else {
				fmt.Println(command + ": not found")
			}
		} else {
			fmt.Println(command, ": command not found")
		}
	}

}
