package internal

import (
	"fmt"
)

func SelectMode() int {
	fmt.Println("Select Mode:")
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Println("3. Exit")

	var choice int
	fmt.Print("Enter choice (1/2/3): ")
	fmt.Scanln(&choice)
	return choice
}

func CallLogin() {
	err := Login("user")
	if err != nil {
		le.Println(err)
	}
}

func CallRegister() {
	err := Register()
	if err != nil {
		le.Println(err)
	}
}
