package session_util

import (
	"fmt"
	"net/http"
)

// CheckAuthHandler checks if the user is logged in
func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, _ := Store.Get(r, "session-name")
	// Check if the user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// If not authenticated, send an error
		fmt.Println("User is NOT authenticated")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("User is authenticated")
	// If authenticated, return success
	w.WriteHeader(http.StatusOK)
}



