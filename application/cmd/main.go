package main

import (
	"chalmers/tkey-group22/application/data/db"
	"chalmers/tkey-group22/application/internal"
	"chalmers/tkey-group22/application/internal/session_util"
	"chalmers/tkey-group22/application/internal/util"
	"fmt"
	"net/http"
)

func main() {

	db, err := db.ConnectMongoDB("mongodb://localhost:27017", "tkeyUserDB")

	// Initialize the User Repository
	internal.UserRepo = util.NewUserRepo(db.Database)

	if err != nil || db == nil || internal.UserRepo == nil {
		fmt.Println("Error connecting to MongoDB")
		return
	}

	http.HandleFunc("/api/register", internal.RegisterHandler)
	http.HandleFunc("/api/login", internal.LoginHandler)
	http.HandleFunc("/api/verify", internal.VerifyHandler)
	http.Handle("/api/initialize-login", enableCors(http.HandlerFunc(internal.InitializeLoginHandler)))
	http.Handle("/api/verify_session", enableCors(http.HandlerFunc(session_util.CheckAuthHandler)))
	http.Handle("/getuser", enableCors(session_util.SessionMiddleware(http.HandlerFunc(internal.GetUserHandler))))

	fmt.Println("Mock application running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
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
