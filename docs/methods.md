```go 
package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) Show() {
	fmt.Println(p.Name)
	fmt.Println(p.Age)
}

func (p *Person) Birthday() {
	fmt.Println("Happy birthday")
	p.Age++
}
func main() {
	person := Person{
		Name: "Abhi",
		Age:  24,
	}
	person.Show()
	person.Birthday()
	person.Show()
}

```
