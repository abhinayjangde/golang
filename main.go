package main

import (
	"fmt"
	"os"
)

func main() {

	data, err := os.ReadFile("README.md")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

}
