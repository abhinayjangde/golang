package main

import "fmt"

func main() {

	// age := 22

	// if age >= 18 {
	// 	fmt.Println("Age is greater than 18")
	// } else {
	// 	fmt.Println("Age is less than 18")
	// }

	// else-if

	// age := 22
	// if age >= 18 {
	// 	fmt.Println("person is adult")
	// } else if age >= 12 {
	// 	fmt.Println("person is teenager")
	// } else {
	// 	fmt.Println("person is a kid")
	// }

	// var role string = "admin"
	// var hasPermissions bool = false

	// if role == "admin" || hasPermissions {
	// 	fmt.Println("yes")
	// }

	// if age := 15; age >= 18 {
	// 	fmt.Println("person is an adult ", age)
	// } else if age >= 12 {
	// 	fmt.Println("person is teenager ", age)
	// }

	// go does not have ternary, you will have to use normal if else.

	var i int = 0;
	for i <= 10 {
		if i == 4 {
			fmt.Println("4 is found")
			break
		}
		// fmt.Println("run", i)
		i = i + 1
	}

	
}
