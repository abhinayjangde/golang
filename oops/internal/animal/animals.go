package animals

import "fmt"

type Animal struct {
	name string
}

func (a Animal) Eat() {
	fmt.Println(a.name, "is eating")
}

// inheritance through composition and embedding
type Dog struct {
	Animal
	breed string
}

func NewAnimal(name string) *Animal {
	return &Animal{
		name: name,
	}
}

func NewDog(name string, breed string) *Dog {
	return &Dog{
		name:  name,
		breed: breed,
	}
}
