package handlers

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

// RegisterHandler handles the user registration process
// It expects a POST request with a JSON body containing the username and public key with label of the user to be registered
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 400 Bad Request: if the request body is invalid or cannot be parsed
// - 409 Conflict: if the user already exists
// - 500 Internal Server Error: if there is an error creating the user or sending the response
// - 200 OK: if the user is registered successfully
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	username := requestBody.Username
	pubkey := requestBody.Pubkey
	label := requestBody.Label

	fmt.Printf("Received registration request for user: %s\n", username)

	// Check if user already exists
	userExists, err := UserRepo.GetUser(username)

	// Checks for sanitization error
	if _, ok := err.(*structs.ErrorInputNotSanitized); ok {
		fmt.Println("Input is not sanitized")
		errMsg := err.(*structs.ErrorInputNotSanitized).Error()
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	if userExists != nil || err != mongo.ErrNoDocuments {
		fmt.Printf("User already exists: %s\n", username)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Store new user data
	user, err := UserRepo.CreateUser(username, pubkey, label)

	// Checks for sanitization error
	if _, ok := err.(*structs.ErrorInputNotSanitized); ok {
		fmt.Println("Input is not sanitized")
		errMsg := err.(*structs.ErrorInputNotSanitized).Error()
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	// Check for other errors
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
