package main

import "fmt"

func main() {
	slice := []int{1, 2, 3}
	// fmt.Println(slice)
	slice = append(slice, 4, 5)
	fmt.Println(slice)

}
