// Package starts the backend server and connects to the MongoDB database
package main

import (
	"chalmers/tkey-group22/application/data/db"
	"chalmers/tkey-group22/application/internal"
	"chalmers/tkey-group22/application/internal/session_util"
	"chalmers/tkey-group22/application/internal/util"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allow frontend
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Needed for cookies/sessions

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Starts the application
func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Connects to the MongoDB dat.env.exampleabase named tkeyUserDB
	db, err := db.ConnectMongoDB(os.Getenv("MONGO_URI"), "tkeyUserDB")

	// Initialize the UserRepository struct with the database reference
	internal.UserRepo = util.NewUserRepo(db.Database)

	if err != nil || db == nil || internal.UserRepo == nil {
		return
	}

	http.HandleFunc("/api/register", internal.RegisterHandler)
	http.Handle("/api/login", enableCORS(http.HandlerFunc(internal.LoginHandler)))
	http.Handle("/api/verify", enableCORS(http.HandlerFunc(internal.VerifyHandler)))
	http.Handle("/api/verify-session", http.HandlerFunc(session_util.CheckAuthHandler))
	http.Handle("/api/getuser", http.HandlerFunc(internal.GetUserHandler))
	http.Handle("/api/add-public-key", (http.HandlerFunc(internal.AddPublicKeyHandler)))
	http.Handle("/api/remove-public-key", http.HandlerFunc(internal.RemovePublicKeyHandler))
	http.Handle("/api/get-public-key-labels", http.HandlerFunc(internal.GetPublicKeyLabelsHandler))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
