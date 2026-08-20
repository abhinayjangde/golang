package main

import "fmt"

func main() {
	user := make(map[string]any)

	user["name"] = "abhi"
	user["age"] = 24
	user["email"] = "abhi@gmail.com"

	fmt.Println(user)

	delete(user, "age")

	fmt.Println(user)

}
