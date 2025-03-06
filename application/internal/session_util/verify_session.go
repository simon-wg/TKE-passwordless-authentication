package session_util

import (
	"net/http"
)

// CheckAuthHandler is an HTTP handler function that checks if a user is authenticated.
// It retrieves the session from the request and verifies the "authenticated" value.
// If the user is not authenticated, it responds with an "Unauthorized" error and a 401 status code.
// If the user is authenticated, it responds with a 200 status code.
func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// This protects routes from being accessed if the user is not logged in
// SessionMiddleware is a middleware function that checks if a user session is authenticated.
// If the session is not authenticated, it responds with an "Unauthorized" status and does not
// call the next handler. If the session is authenticated, it calls the next handler in the chain.
//
// Parameters:
// - next: The next http.Handler to be called if the session is authenticated.
//
// Returns:
// - http.Handler: A new handler that wraps the provided handler with session authentication logic.
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-name")



		auth, ok := session.Values["authenticated"].(bool)

		if !ok || !auth {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r) // Call the next handler
	})
}
