package main

import (
	"fmt"
)

func main() {

	// simple switch
	/*
		i := 2
		switch i {
		case 1:
			fmt.Println("case 1")
		case 2:
			fmt.Println("case 2")
		case 3:
			fmt.Println("case 3")
		default:
			fmt.Println("other")
		}
	*/

	// multiple condition switch
	// switch time.Now().Weekday() {
	// case time.Sunday, time.Saturday:
	// 	fmt.Println("It's weekend.")
	// default:
	// 	fmt.Println("It's workday.")
	// }

	// type switch

	whoAmI := func(i interface{}) {
		switch t := i.(type) {
		case int:
			fmt.Println("It's an integer.")
		case string:
			fmt.Println("It's an string.")
		case bool:
			fmt.Println("It's a boolean.")
		default:
			fmt.Println("Other", t)
		}
	}
	whoAmI("golang")
	// whoAmI(3)
	// whoAmI(true)
	// whoAmI(243.1)

}
