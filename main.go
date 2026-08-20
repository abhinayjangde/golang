package main

import "fmt"

func greet() string {
	// msg = "hello"
	return "hello ji"
}

func add(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func print(values ...any) {
	fmt.Println(values...)
}
func main() {

}
