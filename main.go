package main

import (
	"fmt"
	"sync"
	"time"
)

func f(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("done ", n)

	time.Sleep(time.Millisecond * 50)
}
func main() {
	var wg sync.WaitGroup

	wg.Add((10))
	for i := range 10 {
		f(i, &wg)
	}
	wg.Wait()
}
