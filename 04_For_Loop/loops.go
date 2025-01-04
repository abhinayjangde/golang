package main

import "fmt"

// for - is only looping statement in go
func main() {

	// while loop
	// i := 1
	// for i <= 3 {
	// 	fmt.Println(i)
	// 	i = i + 1
	// }

	// Infinite Loop
	// for {
	// 	fmt.Println("Abhi")
	// }

	// classic for loop
	// for i := 0; i < 3; i++ {
	// 	fmt.Println(i)
	// }

	// break and continue
	// for i := 1; i <= 4; i++ {
	// 	// break
	// 	if i == 2 {
	// 		continue
	// 	}
	// 	fmt.Println(i)
	// }

	//range
	for i := range 3 {
		fmt.Println(i)
	}
}
