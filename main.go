package main

import "fmt"

type Order struct {
	owner  string
	amount int64
}

func (o Order) Owner() string {
	return o.owner
}

func (o *Order) setOwner(owner string) {
	o.owner = owner
}
func main() {
	o := &Order{}

	o.setOwner("abhi")
	fmt.Println(o.Owner())
}
