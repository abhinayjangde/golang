package main

import (
	"context"
	"fmt"
	"time"
)

func processOrders(ctx context.Context, orders <-chan int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker stopping:", ctx.Err())
			return
		case order, ok := <-orders:
			if !ok {
				fmt.Println("no more orders")
				return
			}
			fmt.Println("processed order", order)
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	orders := make(chan int)
	go processOrders(ctx, orders)

	orders <- 101
	orders <- 102

	time.Sleep(10 * time.Millisecond)
	cancel()
	time.Sleep(10 * time.Millisecond)

}
