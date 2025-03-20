package handlers

import (
	"net/http"
)

// This handler returns the username of the current session user
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authenticated user
	username, err := getAuthenticatedUser(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	// Send success response
	response := map[string]string{"message": "Access granted", "user": username}
	sendJSONResponse(w, http.StatusOK, response)

}
