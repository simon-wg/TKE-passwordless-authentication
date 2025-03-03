package session_util

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

func SetSession(w http.ResponseWriter, r *http.Request, username string) error {
    session, err := Store.Get(r, "session-name")
    if err != nil {
        fmt.Println("Error getting session:", err)
        return err
    }


    session.Values["authenticated"] = true
    session.Values["username"] = username

 
    session.Options = &sessions.Options{
        Path:     "/",
        MaxAge:   3600,   
        HttpOnly: true,  
        Secure:   false, 
    }

    err = session.Save(r, w)
    if err != nil {
        fmt.Println("Error saving session:", err)
        return err
    }

    fmt.Println("Session successfully created!")
    return nil
}
