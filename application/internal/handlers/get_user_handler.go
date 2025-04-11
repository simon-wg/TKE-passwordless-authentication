package handlers

import (
	"net/http"
)

// GetUserHandler returns the username of the current session user
// It expects a valid authenticated session
//
// Possible responses:
// - 401 Unauthorized: if the user is not authenticated
// - 200 OK: if the user is authenticated successfully
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user
	username, err := getAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Send success response
	response := map[string]string{"message": "Access granted", "user": username}
	sendJSONResponse(w, http.StatusOK, response)

}
