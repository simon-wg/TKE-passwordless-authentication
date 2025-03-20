package handlers

import (
	"chalmers/tkey-group22/application/internal/structs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetPublicKeyLabelsHandler handles the retrieval of public key labels for a user
// It expects a POST request with a JSON body containing the username
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 400 Bad Request: if the request body is invalid or cannot be parsed
// - 404 Not Found: if the user does not exist
// - 500 Internal Server Error: if there is an error retrieving the labels or sending the response
// - 200 OK: if the labels are retrieved successfully
func GetPublicKeyLabelsHandler(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user
	username, err := getAuthenticatedUser(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("Received request to get public key labels for user: %s\n", username)

	userExists, err := UserRepo.GetUser(username)
	if userExists == nil || err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	labels, err := UserRepo.GetPublicKeyLabels(username)
	if err != nil {
		http.Error(w, "Unable to retrieve public key labels", http.StatusInternalServerError)
		return
	}

	// Send the response
	response := map[string][]string{"labels": labels}
	sendJSONResponse(w, http.StatusOK, response)

}

// AddPublicKeyHandler handles the addition of a new public key for a user
// It expects a POST request with a JSON body containing the username and the new public key
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 400 Bad Request: if the request body is invalid or cannot be parsed
// - 404 Not Found: if the user does not exist
// - 409 Conflict: if the user already has the maximum number of public keys or the label already exists
// - 500 Internal Server Error: if there is an error adding the public key or sending the response
// - 200 OK: if the public key is added successfully
func AddPublicKeyHandler(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user
	username, err := getAuthenticatedUser(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.AddPublicKeyRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newPubKey := requestBody.Pubkey
	label := requestBody.Label

	if label == "" {
		http.Error(w, "Label cannot be empty", http.StatusBadRequest)
	}

	if len(newPubKey) == 0 {
		http.Error(w, "Public key cannot be empty", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received request to add public key for user: %s\n", username)

	userExists, err := UserRepo.GetUser(username)
	if userExists == nil || err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_, err = UserRepo.AddPublicKey(username, newPubKey, label)
	if err != nil {
		if err.Error() == "user already has the maximum number of public keys" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if err.Error() == "public key already exists for the user" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if err.Error() == "label already exists for the user" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Unable to add public key", http.StatusInternalServerError)
		}
		return
	}

	// Send the response
	response := map[string]string{"message": "Public key added successfully"}
	sendJSONResponse(w, http.StatusOK, response)
}

// RemovePublicKeyHandler handles the removal of a public key for a user
// It expects a POST request with a JSON body containing the username and the public key to be removed
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 400 Bad Request: if the request body is invalid or cannot be parsed
// - 404 Not Found: if the user does not exist or the label is not found
// - 409 Conflict: if the user has only one public key
// - 500 Internal Server Error: if there is an error removing the public key or sending the response
// - 200 OK: if the public key is removed successfully
func RemovePublicKeyHandler(w http.ResponseWriter, r *http.Request) {

	// Get the authenticated user
	username, err := getAuthenticatedUser(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.RemovePublicKeyRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	label := requestBody.Label

	if label == "" {
		http.Error(w, "Label cannot be empty", http.StatusBadRequest)
	}

	if username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received request to remove public key for user: %s\n", username)

	userExists, err := UserRepo.GetUser(username)
	if userExists == nil || err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_, err = UserRepo.RemovePublicKey(username, label)
	if err != nil {
		if err.Error() == "user must have at least two public keys to remove one" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else if err.Error() == "specified public key to be removed is not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Unable to remove public key", http.StatusInternalServerError)
		}
		return
	}

	// Send the response
	response := map[string]string{"message": "Public key removed successfully"}
	sendJSONResponse(w, http.StatusOK, response)

}
