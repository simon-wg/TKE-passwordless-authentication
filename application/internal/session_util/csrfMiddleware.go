package session_util

import "github.com/gorilla/csrf"

var CsrfMiddleware = csrf.Protect(
	[]byte("your-32-byte-secret-key-here"),
	csrf.Secure(false), // true in production with HTTPS
	csrf.Path("/"),
	csrf.HttpOnly(true),
	csrf.SameSite(csrf.SameSiteLaxMode),
)
