package session_util

import (
	"fmt"
	"net/http"
)

// CheckAuthHandler checks if the user is logged in
func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Println("User is NOT authenticated")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("User is authenticated")
	w.WriteHeader(http.StatusOK)
}

func SessionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session, _ := Store.Get(r, "session-name")
		fmt.Println("SessionMiddleware: Checking session...")
        auth, ok := session.Values["authenticated"].(bool)
        
        if !ok || !auth {
			fmt.Println("SessionMiddleWare could NOT authenticate")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        fmt.Println("SessionMiddleWare Authenticated")
        next.ServeHTTP(w, r) // Call the next handler
    })
}


