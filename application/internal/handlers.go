package internal

import (
	"bytes"
	"chalmers/tkey-group22/application/internal/session_util"
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

	if label == "" {
		http.Error(w, "Label cannot be empty", http.StatusBadRequest)
	}

	if username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received registration request for user: %s\n", username)

	// Check if user already exists
	userExists, err := UserRepo.GetUser(username)

	if userExists != nil || err != mongo.ErrNoDocuments {
		fmt.Printf("User already exists: %s\n", username)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Store new user data
	user, err := UserRepo.CreateUser(username, pubkey, label)
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

	// Send success response
	w.Write(res)
}

// VerifyHandler handles the verification of a user's signature
// It expects a POST request with a JSON body containing "username" and "signature" fields
//
// Possible responses:
// - 405 Method Not Allowed: if the request method is not POST
// - 400 Bad Request: if the request body is invalid or cannot be parsed
// - 404 Not Found: if the user does not exist
// - 401 Unauthorized: if the signature is invalid
// - 200 OK: if the signature is valid
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	requestBody := structs.VerifyRequest{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the specified user is found
	userExists, err := UserRepo.GetUser(requestBody.Username)

	if userExists == nil || err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if publicKey has an active challenge
	if !HasActiveChallenge(requestBody.Username) {
		http.Error(w, "No active challenge found for the user", http.StatusNotFound)
		return
	}

	// Verify the signed response
	valid, err := VerifySignature(requestBody.Username, requestBody.Signature)
	if !valid {
		fmt.Println(err)
		http.Error(w, "Invalid signature!!!", http.StatusUnauthorized)
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
}

// This handler returns the username of the current session user
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := session_util.Store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	response := map[string]string{"message": "Access granted", "user": username}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// This handler calls SetSession, which creates a session with the user
// func SetSessionHandler(w http.ResponseWriter, r *http.Request){
// 	var requestBody map[string]string
// 	apiKey := r.Header.Get("X-API-Key")

// 	if apiKey != "secret-API-key" {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}
// 	username := requestBody["username"]
// 	session_util.SetSession(w, r, username)
// }

func InitializeLoginHandler(w http.ResponseWriter, r *http.Request) {
	targetURL := "http://localhost:6060/api/login"

	// Read and parse JSON body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil || data["username"] == "" {
		http.Error(w, "Invalid JSON or missing username", http.StatusBadRequest)
		return
	}

	// Forward request to backend
	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to reach backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Run SetSession if backend login is successful
	if resp.StatusCode == http.StatusOK {
		if err := session_util.SetSession(w, r, data["username"]); err != nil {
			http.Error(w, "Failed to set session", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body) // Forward response body to client
}

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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.GetPublicKeyLabelsRequest{}
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

	if username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
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

	responseBody := map[string][]string{"labels": labels}
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Printf("Unable to marshal response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
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

	username := requestBody.Username
	newPubKey := requestBody.Pubkey
	label := requestBody.Label

	if label == "" {
		http.Error(w, "Label cannot be empty", http.StatusBadRequest)
	}

	if username == "" {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
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

	responseBody := map[string]string{"message": "Public key added successfully"}
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Printf("Unable to marshal response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
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

	username := requestBody.Username
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

	responseBody := map[string]string{"message": "Public key removed successfully"}
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		fmt.Printf("Unable to marshal response for user: %s\n", username)
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}
