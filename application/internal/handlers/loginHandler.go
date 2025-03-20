package handlers

import (
	"chalmers/tkey-group22/application/internal"
	"chalmers/tkey-group22/application/internal/structs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LoginHandler handles user login requests
// It expects a POST request with a JSON body containing the username of the user attempting to log in
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 400 Bad Request: if the request body is invalid or cannot be parsed
// - 404 Not Found: if the user does not exist
// - 500 Internal Server Error: if there is an error creating the challenge or sending the response
// - 200 OK: if the challenge is generated successfully
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

	if username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received login request for user: %s\n", username)

	// Check if the specified user is found
	userExists, err := UserRepo.GetUser(username)

	// Checks for sanitization error
	if _, ok := err.(*structs.ErrorInputNotSanitized); ok {
		fmt.Println("Input is not sanitized")
		errMsg := err.(*structs.ErrorInputNotSanitized).Error()
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	if userExists == nil || err != nil {
		fmt.Printf("User not found: %s\n", username)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Generate a challenge using public key
	challenge, _ := internal.GenerateChallenge(username)

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

	// Send success response
	w.Write(res)
}
