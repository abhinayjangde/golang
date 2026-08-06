package main

import (
	"fmt"
	"sync"
	"time"
)

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
}
