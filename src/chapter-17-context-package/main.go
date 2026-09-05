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
func main() {
	bridge := make(chan string)

	go func() {
		bridge <- "hello"
		bridge <- "world"
		close(bridge)
	}()
	consume(bridge)

}
