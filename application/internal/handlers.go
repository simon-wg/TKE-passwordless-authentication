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
	// TODO: Use JSON object instead
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
	// TODO: Use json.Unmarshal
	responseBody := map[string]string{"message": "User registered successfully"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		fmt.Printf("Unable to send response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
	}
}

// LoginHandler handles user login requests.
// It ensures the request is a GET, extracts the username from the request body,
// reads user data from a CSV file, finds the associated public key,
// creates a challenge, and sends the challenge back in the response.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to.
//   - r: The http.Request containing the login request.
//
// Returns:
//   - None
//
// Dependencies:
//   - challenge.go
//   - config.go
//   - csvutil.go
//
// Expected JSON format in request body:
//
//	{
//	  "username": "example_username"
//	}
//
// JSON format in response body:
//
//	{
//	  "challenge": "generated_challenge"
//	}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it is a POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract username from req body
	// TODO: Use json.Unmarshal
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

	// If the found user is the same as the requested user,
	// extract the public key. If public key is empty return error
	pubkey, ok := userData[username]
	if !ok || pubkey == "" {
		fmt.Printf("Public key not found for user: %s\n", username)
		http.Error(w, "Public key not found for specified user", http.StatusNotFound)
		return
	}

	fmt.Printf("Found public key for user %s: %s\n", username, pubkey)

	// Generate a challenge using public key
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
	// Ensure it is a POST request
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	// TODO: Use json.Unmarshal
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := requestBody["username"]
	signature := requestBody["signature"]

	// Read user data
	userData, err := Read(UsersFile)
	if err != nil {
		fmt.Printf("Error reading user data: %v\n", err)
		http.Error(w, "Unable to read user data", http.StatusInternalServerError)
		return
	}

	pubkey, exists := userData[username]
	if !exists {
		fmt.Println("No user named " + username + " exists")
	}

	// Find the username associated with the public key
	// var username string
	// for user, key := range userData {
	// 	if key == pubkey {
	// 		username = user
	// 		break
	// 	}
	// }

	// Check if publicKey has an active challenge
	if !HasActiveChallenge(pubkey) {
		fmt.Println("No active challenge found for the public key")
		http.Error(w, "No active challenge found for the public key", http.StatusNotFound)
		return
	}

	// Verify the signed response
	valid, err := VerifySignature(pubkey, signature)
	if !valid {
		fmt.Println(err)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Send success response
	// TODO: Use json.Unmarshal
	responseBody := map[string]interface{}{
		"message":  "Verification successful",
		"userData": map[string]string{username: pubkey},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseBody); err != nil {
		fmt.Printf("Unable to send response: %v\n", err)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}

	fmt.Println("Verification successful")
}
