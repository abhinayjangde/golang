package main

import "fmt"

func counter() func() int {
	var count int = 0

	return func() int {
		count++
		return count
	}
}

func main() {
	nextNumber := counter()
	println(nextNumber())
	println(nextNumber())
	println(nextNumber())

	// Anonymous function
	func() {
		fmt.Println("This is an anonymous function")
	}()
}
