package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handlers for register, login and verify

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	/*
		1. Ensure it is a POST request
		2. Parse request body
		3. Extract username and pubkey
		4. Check that user is not already registered
		5. Store user data
		6. Send success response
	*/

	// Ensure it is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract username and public key
	username, usernameExists := requestBody["username"]
	publicKey, publicKeyExists := requestBody["publicKey"]

	if !usernameExists || !publicKeyExists {
		http.Error(w, "Username and publicKey are required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received registration request for user: %s\n", username)

	// Read existing user data
	userData, err := Read(UsersFile)
	if err != nil {
		fmt.Printf("Error reading user data: %v\n", err)
		http.Error(w, "Unable to read user data", http.StatusInternalServerError)
		return
	}

	// Check if user already exists
	if _, userExists := userData[username]; userExists {
		fmt.Printf("User already exists: %s\n", username)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Store new user data
	userData[username] = publicKey
	if err := Write(UsersFile, username, publicKey); err != nil {
		fmt.Printf("Error writing user data: %v\n", err)
		http.Error(w, "Unable to save user data", http.StatusInternalServerError)
		return
	}

	// Send success response
	responseBody := map[string]string{"message": "User registered successfully"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		fmt.Printf("Unable to send response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
	}
}
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
	challenge, _ := GenerateChallenge(pubkey)

	// Send the challenge in the response
	responseBody := map[string]string{"challenge": challenge}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		fmt.Printf("Unable to send response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
	}
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Verify request is POST
		Check if challenge is active for the user
		Verify challenge for user
		Send success response
	*/

	// Ensure it is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	publicKey := requestBody["publicKey"]
	signature := requestBody["signature"]

	// Check if publicKey has active signature
	if !HasActiveChallenge(publicKey) {
		http.Error(w, "No active challenge found for the public key", http.StatusNotFound)
		return
	}

	// Verify the signed response
	valid, _ := VerifySignature(publicKey, signature)
	if !valid {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Read user data
	userData, err := Read(UsersFile)
	if err != nil {
		fmt.Printf("Error reading user data: %v\n", err)
		http.Error(w, "Unable to read user data", http.StatusInternalServerError)
		return
	}

	// Send success response
	responseBody := map[string]interface{}{
		"message":  "Verification successful",
		"userData": userData,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		fmt.Printf("Unable to send response: %v\n", err)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
	}
}
