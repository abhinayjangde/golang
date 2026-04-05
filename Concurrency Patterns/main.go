package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(url string, wg *sync.WaitGroup, resultChain chan string) {
	defer wg.Done()
	time.Sleep(time.Millisecond * 50)
	fmt.Println("Image processed: ", url)
	resultChain <- url
}

func main() {
	var wg sync.WaitGroup
	resultChain := make(chan string, 7)

	startTime := time.Now()

	wg.Add(7)
	go worker("image_1.png", &wg, resultChain)
	go worker("image_2.png", &wg, resultChain)
	go worker("image_3.png", &wg, resultChain)
	go worker("image_4.png", &wg, resultChain)
	go worker("image_5.png", &wg, resultChain)
	go worker("image_6.png", &wg, resultChain)
	go worker("image_7.png", &wg, resultChain)
	wg.Wait()
	close(resultChain)

	for result := range resultChain {
		fmt.Printf("received: %s\n", result)
	}

	fmt.Printf("It took %s ms.", time.Since(startTime))
	fmt.Println("done")

}
