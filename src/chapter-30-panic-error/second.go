package main

import (
	"fmt"
	"os"
	"strconv"
)

func canVote(age int) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	if age < 0 {
		panic("age cannot be negative")
	}

	return age >= 18
}

func main() {
	age := os.Args[1]
	intAge, _ := strconv.Atoi(age)
	fmt.Println("Age:", intAge)

	result := canVote(intAge)
	fmt.Println("Can vote:", result)
}
