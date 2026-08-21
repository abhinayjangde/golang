package main

import "fmt"

func main() {
	type User struct {
		Name  string
		Email string
	}

	user := User{
		Name:  "John Doe",
		Email: "john@gmail.com",
	}

	fmt.Println(user.Name)
}
