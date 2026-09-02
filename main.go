package main

import "fmt"

func main() {
	var a, b, c int
	fmt.Print("enter first number: ")
	fmt.Scan(&a)
	fmt.Print("enter first number: ")
	fmt.Scan(&b)

	c = a + b
	fmt.Println("sum = ", c)
}
