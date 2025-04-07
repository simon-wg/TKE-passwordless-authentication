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
		fmt.Printf("Failed to load .env file: %v\n", err)
	}

	// Connects to the MongoDB database named tkeyUserDB
	db, err := db.ConnectMongoDB(os.Getenv("MONGO_URI"), "tkeyUserDB")
	if err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}

	if db == nil {
		fmt.Println("Database connection is nil")
		os.Exit(1)
	}

	// Initialize the UserRepository and NoteRepository struct with the database reference
	handlers.UserRepo = util.NewUserRepo(db.Database)
	handlers.NotesRepo = util.NewNotesRepo(db.Database)

	if handlers.UserRepo == nil {
		fmt.Println("Failed to initialize UserRepository")
		os.Exit(1)
	}

	if handlers.NotesRepo == nil {
		fmt.Println("Failed to initialize NotesRepository")
		os.Exit(1)
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
	http.Handle("/api/logout", session_util.SessionMiddleware(http.HandlerFunc(handlers.LogoutHandler)))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
