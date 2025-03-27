package handlers

import (
	"chalmers/tkey-group22/application/internal/session_util"
	"net/http"
)

// LogoutHandler handles user logout requests
// It expects a POST request with a cookie (automatically included with credentials flag)
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 404 Not Found: if the user doesn't have a session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := session_util.TerminateSession(w, r)
	if err != nil {
		http.Error(w, "No active session found", http.StatusNotFound)
	}

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
