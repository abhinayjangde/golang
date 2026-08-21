package main

import (
	"fmt"
	"os"
)

func main() {

	// data, err := os.ReadFile("README.md")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(data))

	// file, err := os.Create("hello.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// defer file.Close()
	// file.WriteString("test")
	dir, err := os.Open(".")
	if err != nil {
		return
	}
	defer dir.Close()

	fileInfors, err := dir.ReadDir(3)

	for _, f := range fileInfors {
		fmt.Println(f.Name())
	}
}
