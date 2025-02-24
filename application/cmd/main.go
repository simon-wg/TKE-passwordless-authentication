package main

import (
	"chalmers/tkey-group22/application/internal"
	"fmt"
	"net/http"
)

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register", internal.RegisterHandler)
	mux.HandleFunc("/api/login", internal.LoginHandler)
	mux.HandleFunc("/api/verify", internal.VerifyHandler)

	// Apply the CORS middleware
	handler := enableCORS(mux)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
