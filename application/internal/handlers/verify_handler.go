package handlers

import (
	"chalmers/tkey-group22/application/internal"
	"chalmers/tkey-group22/application/internal/session_util"
	"chalmers/tkey-group22/application/internal/structs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
	if !internal.HasActiveChallenge(requestBody.Username) {
		http.Error(w, "No active challenge found for the user", http.StatusNotFound)
		return
	}

	// Verify the signed response
	valid, err := internal.VerifySignature(requestBody.Username, requestBody.Signature, UserRepo)
	if !valid {
		fmt.Println(err)
		http.Error(w, "Invalid signature!!!", http.StatusUnauthorized)
		return
	}

	if err := session_util.SetSession(w, r, requestBody.Username); err != nil {
		http.Error(w, "Failed to set session", http.StatusInternalServerError)
		return
	}

	// We don't expect a response body here, so commenting this out for the while
	// sendJSONResponse(w, http.StatusOK, nil)

}
