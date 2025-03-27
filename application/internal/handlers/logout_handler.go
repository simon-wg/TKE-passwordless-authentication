package handlers

import (
	"chalmers/tkey-group22/application/internal/session_util"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure it is a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve session
	session, _ := session_util.Store.Get(r, "session-name")

	// Invalidate session in the server
	session.Options.MaxAge = -1

	// Save changes to remove session
	session.Save(r, w)

	// Deleted the cookie on the browser
	http.SetCookie(w, &http.Cookie{
		Name:     "session-name",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged out successfully"))
}
