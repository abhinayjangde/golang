package main

import (
	"flag"
	"fmt"
)

func main() {
	// allArgs := os.Args
	// Define flags (Name, Default Value, Description)
	wordPtr := flag.String("word", "foo", "a string value")
	numPtr := flag.Int("n", 42, "an integer value")
	forkPtr := flag.Bool("fork", false, "a bool value")

	// You must call flag.Parse() to execute the parsing
	flag.Parse()

	// Flags are pointers, so dereference them with *
	fmt.Println("word:", *wordPtr)
	fmt.Println("number:", *numPtr)
	fmt.Println("fork:", *forkPtr)

	// Access positional arguments left over after flags
	fmt.Println("Tail arguments:", flag.Args())
}
