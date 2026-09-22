package main

import "fmt"

type person struct {
	FirstName  string
	MiddleName *string
	LastName   string
}

func makePointer[T any](t T) *T {
	return &t
}
func main() {
	x := 10
	y := x
	y = 20

	fmt.Println(x, y)
}
