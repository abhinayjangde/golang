package main

import (
	"fmt"
	"time"
)

type order struct {
	id        string
	amount    float64
	status    string
	createdAt time.Time
}

func main() {
	// var order order
	// order.id = "12345"
	// order.amount = 99.99
	// order.status = "Pending"
	// order.createdAt = time.Now()

	myOrder := order{
		id:        "1",
		amount:    50.23,
		status:    "Pending",
		createdAt: time.Now(),
	}

	fmt.Println(myOrder)

}
