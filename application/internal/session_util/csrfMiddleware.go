package session_util

import (
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

var CsrfMiddleware func(http.Handler) http.Handler

func InitCSRF() {
	csrfKey := os.Getenv("CSRF_KEY")
	if csrfKey == "" {
		panic("CSRF_KEY is not set")
	}

	CsrfMiddleware = csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(true),
		csrf.Path("/"),
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteLaxMode),
	)
}
