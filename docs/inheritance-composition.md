Go doesn't have class inheritance. Instead, we use struct embedding (composition)

```go
package main

import "fmt"

type Animal struct {
	Name  string
	Legs  int
	Sound string
}

// constructor
func NewAnimal(name string, legs int, sound string) *Animal {
	return &Animal{
		Name:  name,
		Legs:  legs,
		Sound: sound,
	}
}
func (a Animal) Speak() {
	fmt.Printf("%s says %s\n", a.Name, a.Sound)
}
func (a Animal) Move() {
	fmt.Printf("%s moves with %d legs\n", a.Name, a.Legs)
}

// Dog embeds Animal (composition)
type Dog struct {
	Animal // Embedding (like inheritance)
	Breed  string
}

func NewDog(name string, legs int, sound string, breed string) *Dog {
	return &Dog{
		Name:  name,
		Legs:  legs,
		Sound: sound,
		Breed: breed,
	}
}

// Override Speak method
func (d Dog) Speak() {
	fmt.Printf("%s, %s barks loudly!\n", d.Name, d.Sound)
}

// Cat embeds Animal
type Cat struct {
	Animal
	IsIndoor bool
}

func NewCat(name string, legs int, sound string, isIndorr bool) *Cat {
	return &Cat{
		Name:     name,
		Legs:     legs,
		Sound:    sound,
		IsIndoor: isIndorr,
	}
}

// Custom method for Cat
func (c Cat) Climb() {
	fmt.Printf("%s is climbing a tree\n", c.Name)
}
func main() {
	animal := NewAnimal("tiger", 4, "oooo")
	animal.Speak()
}
```