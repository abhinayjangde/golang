package main

import "fmt"

func changeNum(num *int) {
	*num = 10
	fmt.Println("Inside changeNum:", *num)
}

func main() {
	// var myNum int = 5
	// changeNum(&myNum)
	// fmt.Println("In main:", myNum)

	// var age int = 23
	// var ptr *int = &age

	// fmt.Println("Value of age:", age)
	// fmt.Println("Address of age:", &age)
	// fmt.Println("Value of ptr (address of age):", ptr)
	// fmt.Println("Value at address stored in ptr:", *ptr)
}
