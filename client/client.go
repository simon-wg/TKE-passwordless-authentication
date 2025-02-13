package main

import (
	"chalmers/tkey-group22/internal"
	"fmt"
	"log"
	"os"
)

var le = log.New(os.Stderr, "Error: ", 0)

func main() {

	// Gets mode from user inputs and runs selected mode. Loops until program is told to exit.
	for {
		mode := select_mode()

		switch mode {
		case 1:
			// Perform register
		case 2:
			// Perform login
			call_login()
		case 3:
			// Stop program
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func select_mode() int {
	fmt.Println("Select Mode:")
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Println("3. Exit")

	var choice int
	fmt.Print("Enter choice (1/2/3): ")
	fmt.Scanln(&choice)
	return choice
}

func call_login() {
	err := internal.Login("user")
	if err != nil {
		le.Println(err)
	}
}
