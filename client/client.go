package main

import (
	"chalmers/tkey-group22/internal"
	"log"
	"os"
)

var le = log.New(os.Stderr, "Error: ", 0)

func main() {
	internal.GetChallengeAndSign("user")
}
