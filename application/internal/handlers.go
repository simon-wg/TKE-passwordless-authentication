package internal

import (
	"chalmers/tkey-group22/application/internal/structs"
	"chalmers/tkey-group22/application/internal/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

var UserRepo util.UserRepository

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
	requestBody := structs.RegisterRequest{}
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

	// Check if user already exists
	userExists, err := UserRepo.GetUser(username)

	fmt.Println(userExists)
	fmt.Println(err)

	if userExists != nil || err != mongo.ErrNoDocuments {
		fmt.Printf("User already exists: %s\n", username)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Store new user data
	user, err := UserRepo.CreateUser(username, pubkey)
	if err != nil || user == nil {
		fmt.Printf("Error creating user: %v\n", err)
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
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
	requestBody := structs.LoginRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// If the username field is empty, return a bad request error
	if requestBody.Username == "" {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	}

	// If the username field has a val, put it in a variable
	username := requestBody.Username

	fmt.Printf("Received login request for user: %s\n", username)

	// Check if the specified user is found
	userExists, err := UserRepo.GetUser(username)

	if userExists == nil || err != nil {
		fmt.Printf("User not found: %s\n", username)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Generate a challenge using public key
	challenge, _ := GenerateChallenge(username)

	// Send the challenge in the response
	response := structs.LoginResponse{
		Challenge: challenge,
	}
	res, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Unable to marshal response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
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
	requestBody := structs.VerifyRequest{}
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

	// Check if the specified user is found
	userExists, err := UserRepo.GetUser(requestBody.Username)

	if userExists == nil || err != nil {
		fmt.Printf("User not found: %s\n", requestBody.Username)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if publicKey has an active challenge
	if !HasActiveChallenge(requestBody.Username) {
		fmt.Println("No active challenge found for the user ")
		http.Error(w, "No active challenge found for the user", http.StatusNotFound)
		return
	}
	print("------------------------------")

	// Verify the signed response
	valid, err := VerifySignature(requestBody.Username, requestBody.Signature)
	if !valid {
		fmt.Println(err)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Send success response
	w.Write([]byte(nil))

	// We don't expect a response body here, so commenting this out for the while
	// responseBody := map[string]interface{}{
	// 	"message":  "Verification successful",
	// 	"userData": map[string]string{requestBody.Username: pubkeyString},
	// }
	// responseBodyBytes, err := json.Marshal(responseBody)
	// if err != nil {
	// 	fmt.Printf("Unable to marshal response: %v\n", err)
	// 	http.Error(w, "Unable to send response", http.StatusInternalServerError)
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(responseBodyBytes)

	fmt.Println("Verification successful")
}
