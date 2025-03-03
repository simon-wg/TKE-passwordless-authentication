package session_util

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// SetSession handles setting the session values for authenticated users


func SetSession(w http.ResponseWriter, r *http.Request, username string) error {
    // Get the session
    session, err := Store.Get(r, "session-name")
    if err != nil {
        fmt.Println("Error getting session:", err)
        return err
    }

    // Set session values
    session.Values["authenticated"] = true
    session.Values["username"] = username

    // Set session options (valid for 1 hour, only over HTTPS)
    session.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   3800,   // Expire the session after an hour
        HttpOnly: true,   // Don't allow JavaScript access
        Secure:   false,  // Set to true in production with HTTPS
    }

    // Save the session
    err = session.Save(r, w)
    if err != nil {
        fmt.Println("Error saving session:", err)
        return err
    }

    fmt.Println("Session successfully created!")
    return nil
}
