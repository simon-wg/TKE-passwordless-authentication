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
	http.HandleFunc("/api/login", internal.LoginHandler)
	http.HandleFunc("/api/verify", internal.VerifyHandler)
	http.Handle("/api/verify-session", enableCors(http.HandlerFunc(session_util.CheckAuthHandler)))
	http.Handle("/api/getuser", enableCors(session_util.SessionMiddleware(http.HandlerFunc(internal.GetUserHandler))))
	http.Handle("/api/add-public-key", enableCors(session_util.SessionMiddleware(http.HandlerFunc(internal.AddPublicKeyHandler))))
	http.Handle("/api/remove-public-key", enableCors(session_util.SessionMiddleware(http.HandlerFunc(internal.RemovePublicKeyHandler))))
	http.Handle("/api/get-public-key-labels", enableCors(session_util.SessionMiddleware(http.HandlerFunc(internal.GetPublicKeyLabelsHandler))))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
