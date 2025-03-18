package internal

import (
	"chalmers/tkey-group22/application/internal/session_util"
	"chalmers/tkey-group22/application/internal/structs"
	"chalmers/tkey-group22/application/internal/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserRepo util.UserRepository
var NotesRepo util.NotesRepository

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

// VerifyHandler handles the verification of a user's signature. If the signiture is valid it
// will set add the user to the session storage and return a cookie in the response.
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

	if err := session_util.SetSession(w, r, requestBody.Username); err != nil {
		http.Error(w, "Failed to set session", http.StatusInternalServerError)
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
	session, _ := session_util.Store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
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
	session, _ := session_util.Store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)

	if !ok || username == "" {
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
	session, _ := session_util.Store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)

	if !ok || username == "" {
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
	session, _ := session_util.Store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)

	if !ok || username == "" {
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

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User unregistered successfully"})
}

// GetNotesHandler handles HTTP GET requests to retrieve notes for a signed-in user.
// It checks if the request method is GET, retrieves the username from the session,
// fetches the notes for the user from the NotesRepo, converts the notes to JSON,
// and sends the JSON response back to the client.
//
// If the request method is not GET, it responds with "Invalid request method" and
// a 405 Method Not Allowed status code.
//
// If there is no user signed in, it responds with "No user signed in" and a 401
// Unauthorized status code.
//
// If there is an error marshalling the notes to JSON, it responds with "Unable to
// marshal notes" and a 500 Internal Server Error status code.
//
// The response content type is set to "application/json".
func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username, err := session_util.GetSessionUsername(r)
	if err != nil {
		http.Error(w, "No user signed in", http.StatusUnauthorized)
		return
	}
	notes, _ := NotesRepo.GetNotes(username)

	// Convert notes to JSON
	responseBodyBytes, err := json.Marshal(notes)
	if err != nil {
		http.Error(w, "Unable to marshal notes", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}

// CreateNoteHandler handles the creation of a new note.
// It expects a POST request with a JSON body containing the note details.
// The request body should be in the format:
//
//	{
//	  "name": "note name",
//	  "note": "note content"
//	}
//
// The function performs the following steps:
// 1. Validates that the request method is POST.
// 2. Reads and unmarshals the request body into a SaveNoteRequest struct.
// 3. Retrieves the username from the session.
// 4. Calls the NotesRepo.CreateNote function to save the note.
// 5. Returns a JSON response with a success message and the ID of the created note.
//
// If any step fails, an appropriate HTTP error response is returned.
//
// Parameters:
// - w: http.ResponseWriter to write the response.
// - r: *http.Request containing the request details.
func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.SaveNoteRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	name := requestBody.Name
	note := requestBody.Note
	username, err := session_util.GetSessionUsername(r)
	if err != nil {
		http.Error(w, "No user signed in", http.StatusUnauthorized)
		return
	}

	result, err := NotesRepo.CreateNote(username, name, note)
	if result == nil || err != nil {
		http.Error(w, "Failed to save notes", http.StatusInternalServerError)
		return
	}

	responseBody := map[string]interface{}{
		"message": "Notes saved successfully",
		"id":      result.InsertedID.(primitive.ObjectID).Hex(),
	}
	responseBodyBytes, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "Unable to send response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}

// UpdateNoteHandler handles the HTTP request for updating a note.
// It expects a POST request with a JSON body containing the note details to be updated.
// The request body should match the structs.UpdateNotesRequest structure.
// The handler performs the following steps:
// 1. Validates the request method is POST.
// 2. Reads and unmarshals the request body into a structs.UpdateNotesRequest object.
// 3. Retrieves the username from the session.
// 4. Fetches the current note entry from the repository using the provided note ID.
// 5. Checks if the current user is the owner of the note.
// 6. Updates the note in the repository with the new details.
// 7. Returns a JSON response indicating the success or failure of the update operation.
//
// If any step fails, an appropriate HTTP error response is returned.
func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.UpdateNotesRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username, _ := session_util.GetSessionUsername(r)
	currentEntry, err := NotesRepo.GetNote(requestBody.ID)
	if err != nil {
		http.Error(w, "Error retrieving entry", http.StatusInternalServerError)
	}

	if username != currentEntry.Username {
		http.Error(w, "User not owner of entry", http.StatusUnauthorized)
		return
	}

	result, err := NotesRepo.UpdateNote(requestBody.ID, username, requestBody.Name, requestBody.Note)
	if result == nil || err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	responseBody := map[string]string{"message": "Note updated successfully"}
	responseBodyBytes, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}

// DeleteNoteHandler handles the deletion of a note.
// It expects a DELETE request with a JSON body containing the note ID to be deleted.
// The handler performs the following steps:
// 1. Verifies that the request method is DELETE.
// 2. Reads and unmarshals the request body into a DeleteNoteRequest struct.
// 3. Retrieves the username from the session.
// 4. Fetches the note entry from the repository using the provided note ID.
// 5. Checks if the authenticated user is the owner of the note.
// 6. Deletes the note from the repository if the user is the owner.
// 7. Returns a success message in JSON format if the note is deleted successfully.
//
// If any of the steps fail, an appropriate HTTP error response is returned.
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	requestBody := structs.DeleteNoteRequest{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username, err := session_util.GetSessionUsername(r)
	if err != nil {
		http.Error(w, "No user signed in", http.StatusUnauthorized)
		return
	}

	currentEntry, err := NotesRepo.GetNote(requestBody.ID)
	if err != nil {
		http.Error(w, "Error retrieving entry", http.StatusInternalServerError)
		return
	}

	if username != currentEntry.Username {
		http.Error(w, "User not owner of entry", http.StatusUnauthorized)
		return
	}

	result, err := NotesRepo.DeleteNote(requestBody.ID)
	if result == nil || err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	responseBody := map[string]string{"message": "Note deleted successfully"}
	responseBodyBytes, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBodyBytes)
}
