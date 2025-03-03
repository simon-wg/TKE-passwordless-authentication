// Package starts the backend server and connects to the MongoDB database
package main

import (
	"chalmers/tkey-group22/application/data/db"
	"chalmers/tkey-group22/application/internal"
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
		fmt.Println("Error connecting to MongoDB")
		return
	}

	http.HandleFunc("/api/register", internal.RegisterHandler)
	http.HandleFunc("/api/login", internal.LoginHandler)
	http.HandleFunc("/api/verify", internal.VerifyHandler)

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
