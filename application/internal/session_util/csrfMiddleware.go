package session_util

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

var CsrfMiddleware func(nextHandler http.Handler) http.Handler

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Printf("Failed to load .env file: %v\n", err)
	}

	csrfKey := os.Getenv("CSRF_KEY")
	if csrfKey == "" {
		fmt.Println("CSRF_KEY is not set in the environment!")
	} else {
		CsrfMiddleware = csrf.Protect(
			[]byte(csrfKey),
			csrf.Secure(true),
			csrf.Path("/"),
			csrf.HttpOnly(true),
			csrf.SameSite(csrf.SameSiteLaxMode),
		)
	}
}
