package main

import (
	"fmt"
	"time"
)

func tast(id int) {
	fmt.Println("doing tast ", id)
}

func main() {
	fmt.Println("start")
	for i := 0; i <= 10; i++ {
		go tast(i)
	}

	time.Sleep(time.Second * 2)
	fmt.Println("end")
}
