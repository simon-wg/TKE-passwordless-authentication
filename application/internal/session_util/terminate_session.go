package session_util

import (
	"net/http"
)

// TerminateSession removes the users session from the sever and terminates the cookie in the browser
//
// Parameters:
//   - w: http.ResponseWriter to write the session cookie to the response.
//   - r: *http.Request to get the session from the request.
//
// Returns:
//   - error: an error if there is an issue getting or saving the session, otherwise nil.
func TerminateSession(w http.ResponseWriter, r *http.Request) error {
	// Retrieve session
	session, err := Store.Get(r, "session-name")
	if err != nil {
		return err
	}

	// Invalidate session in the server
	session.Options.MaxAge = -1

	// Save changes to remove session
	session.Save(r, w)

	// Deletes the cookie on the browser
	http.SetCookie(w, &http.Cookie{
		Name:     "session-name",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	return nil
}
