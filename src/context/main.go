package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	exampleTimeout()
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
