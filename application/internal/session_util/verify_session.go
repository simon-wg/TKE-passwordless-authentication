package session_util

import (
	"net/http"
)

// This protects routes from being accessed if the user is not logged in
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-name")
		_, ok := session.Values["username"]

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r) // Call the next handler
	})
}
