package main

import (
	"chalmers/tkey-group22/internal"
	"log"
	"os"
)

var le = log.New(os.Stderr, "Error: ", 0)

func main() {
	err := internal.Register()
	if err != nil {
		le.Println(err)
	}

	err = internal.Login("user")
	if err != nil {
		le.Println(err)
	}
}
