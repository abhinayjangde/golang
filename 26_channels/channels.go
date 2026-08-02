package main

import (
	"fmt"
	"time"
)

// func processNum(numChan chan int) {

// 	for num := range numChan {
// 		fmt.Println("processing number", num)
// 		time.Sleep(time.Second)
// 	}
// }

//	func sum(result chan int, num1 int, num2 int) {
//		sum := num1 + num2
//		result <- sum
//	}
func main() {
	// numChan := make(chan int)
	// go processNum(numChan)

	// for {
	// 	numChan <- rand.Intn(100)
	// }

	// result := make(chan int)

	// go sum(result, 4, 5)

	// res := <-result

	// fmt.Println(res)

	ch := make(chan int)

	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
			time.Sleep(time.Millisecond * 500)
		}
		close(ch)
	}()

	for val := range ch {
		fmt.Println("received: ", val)
	}
}
