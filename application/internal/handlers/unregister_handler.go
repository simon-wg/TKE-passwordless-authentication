package handlers

import (
	"chalmers/tkey-group22/application/internal/session_util"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

// UnregisterHandler handles user unregistration requests.
// It checks that that request it authorized. Then it ensures the request is a POST, extracts the username from the session.
// checks that the user exists in the database, deletes the user from the database if they exist,
// and then sends a success response.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to.
//   - r: The http.Request containing the unregistration request.
//
// Returns:
//   - None
//
// Dependencies:
//   - UserRepo.go
//
// JSON format in response body:
//
//	{
//	  "message": "User unregistered successfully"
//	}

func UnregisterHandler(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user
	username, err := getAuthenticatedUser(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ensure it is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("Received unregistration request from user: %s\n", username)

	// Check that the user exists in the database
	userExists, err := UserRepo.GetUser(username)
	if userExists == nil || err == mongo.ErrNoDocuments {
		fmt.Printf("User does not exist: %s\n", username)
		http.Error(w, "Could not unregister. User does not exist", http.StatusNotFound)
		return
	}

	// Delete user from the database
	user, err := UserRepo.DeleteUser(username)
	if err != nil || user == nil {
		fmt.Printf("Error deleting user: %v\n", err)
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}

	err = session_util.TerminateSession(w, r)

	if err != nil {
		http.Error(w, "Unable to terminate session", http.StatusInternalServerError)
	}

	// Send success response
	sendJSONResponse(w, http.StatusOK, map[string]string{"message": "User unregistered successfully"})

}
