package main

import "fmt"

// func changeNum(num int) {
// 	num = 5
// 	fmt.Println("in changeNum", num)
// }

func changeNum(num *int) {
	*num = 5
	fmt.Println("in changeNum", *num)
}
func main() {
	num := 1

	changeNum(&num)

	fmt.Println("after changeNum in main function ", num)
}
