package main

import (
	"fmt"
)

type Person struct {
	name     string
	email    string
	password string
}

func (p Person) Details() {
	fmt.Println(p.name)
	fmt.Println(p.email)
}

func (p *Person) SetName(name string) {
	p.name = name
}
func (p *Person) SetEmail(email string) {
	p.email = email
}
func (p *Person) SetPassword(password string) {
	p.password = password
}

func (p Person) GetName() string {
	return p.name
}
func (p Person) GetEmail() string {
	return p.email
}
func main() {
	// inheritance is a (is-a) relationship
	var p Person

	p.SetName("abhi")
	p.SetEmail("abhi@gmail.com")

	// p.Details()
	fmt.Println(p.GetName())
}
