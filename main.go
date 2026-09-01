package main

import "fmt"

func main() {
	var a, b, c byte

	fmt.Scan(&a)
	fmt.Scan(&b)

	c = a + b
	fmt.Println(c)
}
