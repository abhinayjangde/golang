package main

import (
	"context"
	"fmt"
)

func main() {
	root := context.Background()
	fmt.Println("background deadline:", deadlineString(root))
	fmt.Println("background err:", root.Err())

	placeholder := context.TODO()
	fmt.Println("todo deadline:    ", deadlineString(placeholder))
	fmt.Println("todo err:         ", placeholder.Err())
}

func deadlineString(ctx context.Context) string {
	d, ok := ctx.Deadline()
	if !ok {
		return "none"
	}
	return d.String()
}
