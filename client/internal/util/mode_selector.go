package util

import (
	"bufio"
	"chalmers/tkey-group22/internal/auth"
	"fmt"
	"log"
	"os"
	"strings"
)

var le = log.New(os.Stderr, "Error: ", 0)
var appurl = "http://localhost:8080"

func SelectMode() int {
	fmt.Println("\nSelect Mode:")
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Println("3. Exit")

	var choice int
	fmt.Print("Enter choice (1/2/3): ")
	fmt.Scanln(&choice)
	return choice
}

func CallLogin() {
	username := getUsername()
	err := auth.Login(appurl, username)
	if err != nil {
		le.Println(err)
	}
}

func CallRegister() {
	username := getUsername()
	err := auth.Register(appurl, username)
	if err != nil {
		le.Println(err)
	}
}

func getUsername() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter username: ")
	username, _ := reader.ReadString('\n')
	return strings.TrimSpace(username)
}
