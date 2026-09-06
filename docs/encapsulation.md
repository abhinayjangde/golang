```go
package main

import "fmt"

// encapsulation

type BankAccount struct {
	owner   string
	balance float64
}

func (b *BankAccount) deposit(amount float64) {
	if amount > 0 {
		b.balance += amount
		fmt.Println("money deposit to your account")
	}
}

func (b *BankAccount) withdraw(amount float64) {
	if b.balance >= amount {
		b.balance -= amount
		fmt.Println("money withdraw from your account")
	}
}

func (b BankAccount) getBalance() {
	fmt.Printf("Your current balance is %f", b.balance)
}
func main() {
	abhi := BankAccount{
		owner:   "Abhi",
		balance: 1000,
	}

	abhi.withdraw(1000)
	abhi.deposit(200)
	abhi.getBalance()
}

```