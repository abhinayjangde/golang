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
	add := func(x, y int) int {
		return x + y
	}
	x := 0
	increment := func() int {
		x++
		return x
	}
	print(add(1, 2))
	print(increment())
	print(increment())
	print(increment())
}
