package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handlers for register, login and verify
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	/*
		1. Check if valid request
		2. Extract username from request body
		3. Read user data from csv
		4. Find the associated public key
		5. Create a challenge
		6. Send the challenge back in response

	*/
	// Ensure it is a GET
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract username from req body
	var requestBody map[string]string

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// If the username field has a val, put it in a variable
	username, ok := requestBody["username"]
	if !ok {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received login request for user: %s\n", username)

	// Read user data from csv and store in var userData
	userData, err := Read(UsersFile)
	if err != nil {
		fmt.Printf("Error reading user data: %v\n", err)
		http.Error(w, "Unable to read user data", http.StatusInternalServerError)
		return
	}

	// Check if the specified user is found
	if _, userExists := userData[username]; !userExists {
		fmt.Printf("User not found: %s\n", username)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// If username == username that is stored in userData, extract public key in
	// var pubkey. If public key is empty return error
	pubkey, ok := userData[username]
	if !ok || pubkey == "" {
		fmt.Printf("Public key not found for user: %s\n", username)
		http.Error(w, "Public key not found for specified user", http.StatusNotFound)
		return
	}

	fmt.Printf("Found public key for user %s: %s\n", username, pubkey)

	// Generate a challenge using username
	challenge := "challenge"

	// Send the challenge in the response
	responseBody := map[string]string{"challenge": challenge}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		fmt.Printf("Unable to send response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
	}
}
