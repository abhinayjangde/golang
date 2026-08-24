package main

import (
	animals "github.com/abhinayjangde/oops/internal/animal"
	interfaces "github.com/abhinayjangde/oops/internal/interface"
	"github.com/abhinayjangde/oops/internal/person"
)

func main() {
	p := person.NewPerson("abhi", "abhi@gmail.com")
	p.SetName("aditi")
	// fmt.Println(p.GetName())
	// inheritance
	d := animals.NewDog("Romi", "indian")
	d.Eat()

	card := interfaces.CreditCard{}
	upi := interfaces.UPI{}

	interfaces.ProcessPayment(card, 1000)
	interfaces.ProcessPayment(upi, 500)
}
