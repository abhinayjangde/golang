package main

import "fmt"

type Recipient struct {
	Name  string // public field if you use uppercase, if lowercase it is private
	Email string // public field if you use uppercase, if lowercase it is private
}

func main() {
	fmt.Println("Hello, Email Dispatcher!")
	err := loadRecipients("./emails.csv")
	if err != nil {
		fmt.Println("Error loading recipients:", err)
	}
}
