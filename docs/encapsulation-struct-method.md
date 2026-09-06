Go takes a unique approach to OOP. Instead of traditional classes and inheritance, Go uses structs and interfaces with composition over inheritance. Let me teach you the core OOP concepts in Go!

```go
package main

import "fmt"

// this is like a - class
type Person struct {
	name string // private field (lowercase) - only accessible within the package
	age  int

	Email string // public field (start with uppercase)
}

// constructor function (convention: NewStructName)

func NewPerson(name string, age int, email string) *Person {
	return &Person{
		name:  name,
		age:   age,
		Email: email,
	}
}

func (p Person) GetName() string {
	return p.name
}

func (p *Person) SetAge(age int) {
	p.age = age
}

func (p Person) GetAge() int {
	return p.age
}
func (p Person) IsAdult() bool {
	return p.age >= 18
}

func (p Person) DisplayInfo() {
	fmt.Printf("Name: %s, Age: %d, Email: %s\n", p.name, p.age, p.Email)
}
func main() {
	p := NewPerson("abhi", 24, "abhi@gmail.com")
	p.SetAge(21)
	p.DisplayInfo()

	fmt.Println(p.GetName(), "is adult:", p.IsAdult())
}
```