package session_util

import (
	"net/http"
)

// SessionMiddleware is a middleware function that checks for the existence of a session
// and verifies if the "username" key is present in the session values. If the "username"
// key is not found, it responds with an "Unauthorized" error and a 401 status code.
// If the "username" key is found, it calls the next handler in the chain.
//
// Parameters:
// - next: The next http.Handler to be called if the session is valid.
//
// Returns:
// - http.Handler: A handler that wraps the provided handler with session validation logic.
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-name")
		_, ok := session.Values["username"]

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
