package main

import (
	"fmt"
	"strings"
)

type User struct {
	Name  string
	Email string
}

func (u *User) getEmail() string {
	return u.Email
}
func main() {

	uname := "abhinay jangde"
	fmt.Println(strings.Split(uname, " "))
	for item := range strings.SplitSeq(uname, " ") {
		fmt.Println(item)
	}
}
