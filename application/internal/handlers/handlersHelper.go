package handlers

import (
	"chalmers/tkey-group22/application/internal/session_util"
	"encoding/json"
	"fmt"
	"net/http"
)

// Helper function to retrieve the authenticated username from session
func getAuthenticatedUser(r *http.Request) (string, error) {
	session, _ := session_util.Store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		return "", fmt.Errorf("unauthorized")
	}
	return username, nil
}

// Helper function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
