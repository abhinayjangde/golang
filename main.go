package main

import (
	"fmt"
	"time"
)

type Product struct {
	ID          int
	Name        string
	Price       float64
	Stock       int
	Description string
	CreatedAt   time.Time
}

func (p *Product) IsAvailable(quantity int) bool {
	return p.Stock >= quantity
}

func (p *Product) ReduceStock(quantity int) error {
	if !p.IsAvailable(quantity) {
		return fmt.Errorf("insuffiecient stock: available %d, requested %d", p.Stock, quantity)
	}
	p.Stock -= quantity
	return nil
}

type InventoryManager interface {
	CheckStock(productID int, quantity int) bool
	ReserveProduct(productID int, quantity int) error
	ReleaseProduct(productID int, quantity int)
}

type PaymentProcessor interface {
	ProcessPayment(amount float64, method string) (string, error)
	Refund(transactionID string) error
}

type NotificationSender interface {
	SendOrderConfirmation(orderID int, email string) error
	SendOrderStatusUpdate(orderID int, status string, email string) error
}

func main() {
	p := Product{}
}
