package session_util

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// SetSession creates a new session for the given username and saves it in the session store.
// It sets various session options such as path, max age, HttpOnly, Secure, and SameSite mode.
//
// Parameters:
//   - w: http.ResponseWriter to write the session cookie to the response.
//   - r: *http.Request to get the session from the request.
//   - username: string representing the username to be stored in the session.
//
// Returns:
//   - error: an error if there is an issue getting or saving the session, otherwise nil.

func SetSession(w http.ResponseWriter, r *http.Request, username string) error {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		fmt.Println("Error getting session:", err)
		return err
	}

	session.Values["username"] = username

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	if err := session.Save(r, w); err != nil {
		fmt.Println("Error saving session:", err)
		return err
	}

	return nil
}
