package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := fetchData(ctx)

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func fetchData(ctx context.Context) error {
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("Data fetched successfully")
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
func exampleTimeout() {
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*4)
	defer cancel()

	done := make(chan string)

	go func() {
		time.Sleep(time.Second * 3)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("called The API")
	case <-ctxWithTimeout.Done():
		fmt.Println("timeout expired:", ctxWithTimeout.Err())
	}
}
