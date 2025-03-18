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
	internal.PasswordRepo = util.NewPasswordRepo(db.Database)

	if err != nil || db == nil || internal.UserRepo == nil {
		return
	}

	http.HandleFunc("/api/register", internal.RegisterHandler)
	http.Handle("/api/login", http.HandlerFunc(internal.LoginHandler))
	http.Handle("/api/verify", http.HandlerFunc(internal.VerifyHandler))
	http.Handle("/api/getuser", http.HandlerFunc(internal.GetUserHandler))
	http.Handle("/api/unregister", http.HandlerFunc(internal.UnregisterHandler))
	http.Handle("/api/add-public-key", (http.HandlerFunc(internal.AddPublicKeyHandler)))
	http.Handle("/api/remove-public-key", http.HandlerFunc(internal.RemovePublicKeyHandler))
	http.Handle("/api/get-public-key-labels", http.HandlerFunc(internal.GetPublicKeyLabelsHandler))
	http.Handle("/api/save-password", session_util.SessionMiddleware(http.HandlerFunc(internal.CreatePasswordHandler)))
	http.Handle("/api/get-user-passwords", session_util.SessionMiddleware(http.HandlerFunc(internal.GetUserPasswordsHandler)))
	http.Handle("/api/update-password", session_util.SessionMiddleware(http.HandlerFunc(internal.UpdatePasswordHandler)))
	http.Handle("/api/delete-password", session_util.SessionMiddleware(http.HandlerFunc(internal.DeletePasswordHandler)))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
