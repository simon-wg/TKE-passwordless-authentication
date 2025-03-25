// Package starts the backend server and connects to the MongoDB database
package main

import (
	"chalmers/tkey-group22/application/data/db"
	"chalmers/tkey-group22/application/internal/handlers"
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

	// Initialize the UserRepository and NoteRepository struct with the database reference
	handlers.UserRepo = util.NewUserRepo(db.Database)
	handlers.NotesRepo = util.NewNotesRepo(db.Database)

	if err != nil || db == nil || handlers.UserRepo == nil {
		return
	}

	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.Handle("/api/login", http.HandlerFunc(handlers.LoginHandler))
	http.Handle("/api/verify", http.HandlerFunc(handlers.VerifyHandler))
	http.Handle("/api/getuser", http.HandlerFunc(handlers.GetUserHandler))
	http.Handle("/api/unregister", http.HandlerFunc(handlers.UnregisterHandler))
	http.Handle("/api/add-public-key", http.HandlerFunc(handlers.AddPublicKeyHandler))
	http.Handle("/api/remove-public-key", http.HandlerFunc(handlers.RemovePublicKeyHandler))
	http.Handle("/api/get-public-key-labels", http.HandlerFunc(handlers.GetPublicKeyLabelsHandler))
	http.Handle("/api/create-note", session_util.SessionMiddleware(http.HandlerFunc(handlers.CreateNoteHandler)))
	http.Handle("/api/get-user-note", session_util.SessionMiddleware(http.HandlerFunc(handlers.GetNotesHandler)))
	http.Handle("/api/update-note", session_util.SessionMiddleware(http.HandlerFunc(handlers.UpdateNoteHandler)))
	http.Handle("/api/delete-note", session_util.SessionMiddleware(http.HandlerFunc(handlers.DeleteNoteHandler)))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
