package main

import (
	"fmt"
	"sync"
)

func tast(id int, w *sync.WaitGroup) {
	defer w.Done()
	fmt.Println("done ", id)
	fmt.Println("doing tast ", id)
}

func main() {
	fmt.Println("start")

	var wg sync.WaitGroup

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go tast(i, &wg)
	}
	wg.Wait()

	fmt.Println("end")
}
