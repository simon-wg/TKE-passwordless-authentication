package session_util

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitSession() {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		fmt.Println("SESSION_KEY is not set in the environment")
	} else {
		Store = sessions.NewCookieStore([]byte(sessionKey))
	}
}
