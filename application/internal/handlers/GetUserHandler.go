package handlers

import (
	"chalmers/tkey-group22/application/internal/session_util"
	"encoding/json"
	"net/http"
)

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
