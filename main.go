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

func main() {
	ch := make(chan int)

	go func() {

		ch <- 5
	}()
	res := <-ch

	fmt.Print(res)
}
