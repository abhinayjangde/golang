package main

import (
	"fmt"
	"sync"
	"time"
)

func print10(n int) int {
	if n == 11 {
		return 0
	}
	fmt.Println(n)
	return print10(n + 1)
}
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("worker", id, "started")
	time.Sleep(time.Second)
	fmt.Println("worker", id, "done")
}

func nums(ch chan<- int) {
	for i := range 5 {
		ch <- i
	}
	close(ch)
}

func main() {
	ch := make(chan int, 5)

	nums(ch)
	for n := range ch {
		fmt.Println(n)
	}

	print10(1)
}
