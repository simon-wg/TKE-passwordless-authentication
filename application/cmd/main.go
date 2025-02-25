package main

import (
	"chalmers/tkey-group22/application/data/db"
	"chalmers/tkey-group22/application/internal"

	"fmt"
	"net/http"
)

func main() {

	db, err := db.ConnectMongoDB("mongodb://localhost:27017", "tkeyUserDB", "users")

	if err != nil || db == nil {
		fmt.Println("Error connecting to MongoDB")
		return
	}

	http.HandleFunc("/api/register", internal.RegisterHandler)
	http.HandleFunc("/api/login", internal.LoginHandler)
	http.HandleFunc("/api/verify", internal.VerifyHandler)

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}
