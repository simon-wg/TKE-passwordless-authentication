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

	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.Handle("/api/login", http.HandlerFunc(handlers.LoginHandler))
	mux.Handle("/api/verify", http.HandlerFunc(handlers.VerifyHandler))
	mux.Handle("/api/getuser", session_util.SessionMiddleware(session_util.CsrfMiddleware((http.HandlerFunc(handlers.GetUserHandler)))))
	mux.Handle("/api/unregister", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.UnregisterHandler))))
	mux.Handle("/api/add-public-key", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.AddPublicKeyHandler))))
	mux.Handle("/api/remove-public-key", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.RemovePublicKeyHandler))))
	mux.Handle("/api/get-public-key-labels", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.GetPublicKeyLabelsHandler))))

	mux.Handle("/api/csrf-token", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.GetCSRF))))

	mux.Handle("/api/create-note", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.CreateNoteHandler))))
	mux.Handle("/api/get-user-note", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.GetNotesHandler))))
	mux.Handle("/api/update-note", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.UpdateNoteHandler))))
	mux.Handle("/api/delete-note", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.DeleteNoteHandler))))
	mux.Handle("/api/logout", session_util.SessionMiddleware(session_util.CsrfMiddleware(http.HandlerFunc(handlers.LogoutHandler))))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
