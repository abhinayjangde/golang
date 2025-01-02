package main

import "fmt"

func main() {

	var name string = "Abhi"
	var age int = 22
	var isActive = true // infer
	var pi float64 = 3.14

	// shorthand
	email := "abhinay@golang.dev"

	var username string

	fmt.Println(name)
	fmt.Println(age)
	fmt.Println(isActive)
	fmt.Println(pi)
	fmt.Println(email)

	username = "golang"
	fmt.Println(username)
}
