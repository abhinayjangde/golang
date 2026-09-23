package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

		command = strings.TrimSpace(command)

		builtins := []string{"exit", "echo", "type"}
		executables := map[string]string{
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
			cmd := exec.Command(executables["git"], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			if err := cmd.Run(); err != nil {
				fmt.Fprintln(os.Stderr, "git error:", err)
			}
		} else if strings.HasPrefix(command, "docker") {
			args := strings.Fields(command)
			cmd := exec.Command(executables["docker"], args[1:]...)
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

			if slices.Contains(builtins, command) {
				fmt.Println(command + " is a shell builtin")
			} else if executables[command] != "" {
				path, err := exec.LookPath(executables[command])

				if err != nil {
					fmt.Fprintln(os.Stderr, "Error looking up command:", err)
					os.Exit(1)
				}
				fmt.Println(command+" is ", path)

			} else {
				fmt.Println(command + ": not found")
			}
		} else {
			fmt.Println(command, ": command not found")
		}
	}

}
