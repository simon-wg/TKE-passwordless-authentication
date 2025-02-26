package main

import (
	"chalmers/tkey-group22/internal/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.Handle("/api/register", enableCors(http.HandlerFunc(registerHandler)))
	http.Handle("/api/login", enableCors(http.HandlerFunc(loginHandler)))

	fmt.Println("Client running on http://localhost:6060")
	http.ListenAndServe(":6060", nil)
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	origin = replaceOriginPort(origin)

	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	username := requestBody["username"]
	err := auth.Login(origin, username)
	if err != nil {
		http.Error(w, "Failed to log in", http.StatusBadRequest)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	origin = replaceOriginPort(origin)

	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	username := requestBody["username"]
	err := auth.Register(origin, username)
	if err != nil {
		http.Error(w, "Failed to register", http.StatusBadRequest)
	}
}

// Change port of request to 8080
func replaceOriginPort(origin string) string {
	parts := strings.Split(origin, ":")
	if len(parts) > 1 {
		parts[len(parts)-1] = "8080"
		origin = strings.Join(parts, ":")
	}
	return origin
}
