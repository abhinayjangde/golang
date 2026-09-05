package main

import (
	"context"
	"fmt"
)

func processOrders(ctx context.Context, orders <-chan int) {

}

func consume(ch <-chan string) {
	for msg := range ch {
		fmt.Println(msg)
	}

}

func streamNumber() <-chan int {
	ch := make(chan int)

	go func() {
		for i := range 4 {
			ch <- i
		}
		close(ch)
	}()
	return ch
}
func main() {
	bridge := make(chan string)

	go func() {
		bridge <- "hello"
		bridge <- "world"
		close(bridge)
	}()
	// consume(bridge)

	dataStream := streamNumber()

	for n := range dataStream {
		fmt.Println(n)
	}

}
