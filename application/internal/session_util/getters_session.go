package session_util

import (
	"fmt"
	"net/http"
)

// Get the username field from the session
func GetSessionUsername(r *http.Request) (string, error) {
	session, _ := Store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in session")
	} else {
		return username, nil
	}
}
