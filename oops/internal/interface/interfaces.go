package interfaces

import "fmt"

type Payment interface {
	Pay(amount float64)
}

type CreditCard struct{}

func (c CreditCard) Pay(amount float64) {
	fmt.Printf("Paid %.2f using Credit Card\n", amount)
}

type UPI struct{}

func (u UPI) Pay(amount float64) {
	fmt.Printf("Paid %.2f using UPI\n", amount)
}

func ProcessPayment(p Payment, amount float64) {
	p.Pay(amount)
}
