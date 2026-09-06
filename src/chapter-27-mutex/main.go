package main

import "fmt"

type post struct {
	views int
}

func (p *post) inc() {
	p.views += 1
}

func main() {
	myPost := post{views: 0}

	myPost.inc()
	myPost.inc()

	fmt.Println(myPost.views)
}
