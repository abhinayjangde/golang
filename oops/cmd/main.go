package main

import (
	animals "github.com/abhinayjangde/oops/internal/animal"
	"github.com/abhinayjangde/oops/internal/person"
)

func main() {
	p := person.NewPerson("abhi", "abhi@gmail.com")
	p.SetName("aditi")
	// fmt.Println(p.GetName())
	// inheritance
	d := animals.NewDog("Romi", "indian")
	d.Eat()
}
