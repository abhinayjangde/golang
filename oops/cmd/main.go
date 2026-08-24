package main

import (
	"fmt"

	"github.com/abhinayjangde/oops/internal/person"
)

func main() {
	p := person.NewPerson("abhi", "abhi@gmail.com")
	p.SetName("aditi")
	fmt.Println(p.GetName())
}
