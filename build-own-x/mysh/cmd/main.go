package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your name:")

	if scanner.Scan() {
		input := scanner.Text()
		fmt.Printf("You entered: %s\n", input)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input: ", err)
	}
}
