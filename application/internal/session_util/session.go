package session_util

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var Store *sessions.CookieStore

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("Failed to load .env file: %v\n", err)
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		fmt.Println("SESSION_KEY is not set in the environment")
	} else {
		Store = sessions.NewCookieStore([]byte(sessionKey))
	}
}
