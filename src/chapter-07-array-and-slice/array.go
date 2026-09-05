package main

import "fmt"

func main() {
	// var nums [4]int
	// nums[0] = 99
	// arr := [5]int{4, 5, 6, 7, 8}
	// fmt.Println(nums)
	// fmt.Println(len(arr))

	// var mark1 = [...]int{12, 32, 54}               // infered size
	// names := [...]string{"Abhi", "Arya", "Chinki"} // infered size
	// fmt.Println(mark1)
	// fmt.Println(names)
	// fmt.Println(names[0])

	// Initialize Only Specific Elements

	// var num1 = [5]int{1: 45, 3: 90}
	// fmt.Println(num1)

	var x string = "Abhi"
	var y int = -12
	fmt.Printf("type %T and value %v\n", x, x)
	fmt.Printf("type %T and value %v\n", y, y)

}
