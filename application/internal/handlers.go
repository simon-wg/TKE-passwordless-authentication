package internal

import (
	"chalmers/tkey-group22/application/internal/util"
	"encoding/json"
	"fmt"
	"io"
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
	requestBody := RegisterRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract username and public key
	username := requestBody.Username
	pubkey := requestBody.Pubkey

	fmt.Printf("Received registration request for user: %s\n", username)

	// Read existing user data
	userData, err := util.Read(UsersFile)
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
	userData[username] = string(pubkey)
	if err := util.Write(UsersFile, username, string(pubkey)); err != nil {
		fmt.Printf("Error writing user data: %v\n", err)
		http.Error(w, "Unable to save user data", http.StatusInternalServerError)
		return
	}

	// Send success response
	responseBody := map[string]string{"message": "User registered successfully"}
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Printf("Unable to marshal response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
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

	// Parse request body
	requestBody := LoginRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// If the username field has a val, put it in a variable
	username := requestBody.Username

	fmt.Printf("Received login request for user: %s\n", username)

	// Read user data from csv and store in var userData
	userData, err := util.Read(UsersFile)
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
	pubkeyString, ok := userData[username]
	if !ok || pubkeyString == "" {
		fmt.Printf("Public key not found for user: %s\n", username)
		http.Error(w, "Public key not found for specified user", http.StatusNotFound)
		return
	}

	fmt.Printf("Found public key for user %s: %s\n", username, pubkeyString)

	// Generate a challenge using public key
	challenge, _ := GenerateChallenge(pubkeyString)

	// Send the challenge in the response
	responseBody := map[string]string{"challenge": challenge}
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Printf("Unable to marshal response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}

// VerifyHandler handles the verification of a user's signature.
// It expects a POST request with a JSON body containing "username" and "signature" fields.
// The handler performs the following steps:
// Request Body:
//
//	{
//	  "username": "exampleUser",
//	  "signature": "hexEncodedSignature"
//	}
//
// Response Body (on success):
//
//	{
//	  "message": "Verification successful",
//	  "userData": {
//	    "exampleUser": "publicKeyString"
//	  }
//	}
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it is a POST request
	if r.Method != http.MethodPost {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	requestBody := VerifyRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Println("Invalid request body")
		http.Error(w, "Invalid signature format", http.StatusBadRequest)
	}

	// Read user data
	userData, err := util.Read(UsersFile)
	if err != nil {
		fmt.Printf("Error reading user data: %v\n", err)
		http.Error(w, "Unable to read user data", http.StatusInternalServerError)
		return
	}

	pubkeyString, exists := userData[requestBody.Username]
	if !exists {
		fmt.Println("No user named " + requestBody.Username + " exists")
	}
	// Check if publicKey has an active challenge
	if !HasActiveChallenge(pubkeyString) {
		fmt.Println("No active challenge found for the public key")
		http.Error(w, "No active challenge found for the public key", http.StatusNotFound)
		return
	}

	// Check if publicKey has an active challenge
	if !HasActiveChallenge(pubkeyString) {
		fmt.Println("No active challenge found for the public key")
		http.Error(w, "No active challenge found for the public key", http.StatusNotFound)
		return
	}

	// Verify the signed response
	valid, err := VerifySignature(pubkeyString, requestBody.Signature)
	if !valid {
		fmt.Println(err)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Send success response
	responseBody := map[string]interface{}{
		"message":  "Verification successful",
		"userData": map[string]string{requestBody.Username: pubkeyString},
	}
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Printf("Unable to marshal response: %v\n", err)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)

	fmt.Println("Verification successful")
}

type VerifyRequest struct {
	Username  string `json:"username"`
	Signature []byte `json:"signature"`
}

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Pubkey   []byte `json:"pubkey"`
}
