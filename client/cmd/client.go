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
		mode := internal.SelectMode()

		switch mode {
		case 1:
			// Perform register
			internal.CallRegister()
		case 2:
			// Perform login
			internal.CallLogin()
		case 3:
			// Stop program
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
