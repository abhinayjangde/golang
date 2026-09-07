package main

import "fmt"

type Order struct {
	owner string
	total string
}

func (o Order) Owner() string {
	return o.owner
}

func (o Order) Total() string {
	return o.total
}

func (o *Order) SetOwner(name string) {
	o.owner = name
}
func main() {
	o := &Order{}
	o.SetOwner("abhi")
	fmt.Println(o.Owner())
}
