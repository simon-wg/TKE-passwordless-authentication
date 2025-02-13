package main

import (
	"chalmers/tkey-group22/internal"
	"fmt"
	"log"
	"os"
)

var le = log.New(os.Stderr, "Error: ", 0)

func main() {

	for {
		mode := select_mode()

		if mode == 1 {
			// Perform register

		} else if mode == 2 {

			// Perform Login
			call_login()

		} else if mode == 3 {
			// Stop program
			break

		} else {
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
