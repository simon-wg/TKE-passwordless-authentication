// Package starts the backend server and connects to the MongoDB database
package main

import (
	"chalmers/tkey-group22/application/data/db"
	"chalmers/tkey-group22/application/internal"
	"chalmers/tkey-group22/application/internal/session_util"
	"chalmers/tkey-group22/application/internal/util"
	"fmt"
	"net/http"
)

// Starts the application
func main() {

	// Connects to the MongoDB database named tkeyUserDB
	db, err := db.ConnectMongoDB("mongodb://localhost:27017", "tkeyUserDB")

	// Initialize the UserRepository struct with the database reference
	internal.UserRepo = util.NewUserRepo(db.Database)

	if err != nil || db == nil || internal.UserRepo == nil {
		return
	}

	http.HandleFunc("/api/register", internal.RegisterHandler)
	http.HandleFunc("/api/login", internal.LoginHandler)
	http.HandleFunc("/api/verify", internal.VerifyHandler)
	http.Handle("/api/initialize-login", enableCors(http.HandlerFunc(internal.InitializeLoginHandler)))
	http.Handle("/api/verify-session", enableCors(http.HandlerFunc(session_util.CheckAuthHandler)))
	http.Handle("/api/getuser", enableCors(session_util.SessionMiddleware(http.HandlerFunc(internal.GetUserHandler))))
	http.Handle("/api/add-public-key", enableCors(session_util.SessionMiddleware(http.HandlerFunc(internal.AddPublicKeyHandler))))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:8080" || origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
