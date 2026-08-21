package main

import (
	"errors"
	"fmt"
)

func main() {
	err := errors.New("my error")
	fmt.Println(err)
}
