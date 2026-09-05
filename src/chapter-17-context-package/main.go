package main

import (
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println(ctx.Err())
	defer cancel()
}
