package main

import (
	"fmt"

	"github.com/abhinayjangde/auth/auth"
	"github.com/abhinayjangde/auth/user"
	"github.com/fatih/color"
)

func main() {
	// The main function is intentionally left empty.
	auth.Login("abhinay", "abhi1234")

	session := auth.GetSession()
	println("Session status:", session)

	user := user.User{
		Email: "abhinay@example.com",
		Name:  "Abhinay Jangde",
	}

	fmt.Println(user.Email, user.Name)
	color.Cyan("Prints text in cyan.")
	color.Blue("Prints %s in blue.", "text")
	color.Red("We have red")
	color.RGB(255, 128, 0).Println("foreground orange")
	color.RGB(230, 42, 42).Println("foreground red")

	color.BgRGB(255, 128, 0).Println("background orange")
	color.BgRGB(230, 42, 42).Println("background red")
}
