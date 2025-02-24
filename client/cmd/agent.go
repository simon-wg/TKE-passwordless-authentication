package main

// // import (
// // 	"chalmers/tkey-group22/internal/tkey"
// // 	"encoding/json"
// // 	"fmt"
// // 	"io/ioutil"
// // 	"net/http"
// // )

// // func enableCORS(next http.Handler) http.Handler {
// // 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// // 		w.Header().Set("Access-Control-Allow-Origin", "*")
// // 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// // 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// // 		if r.Method == "OPTIONS" {
// // 			return
// // 		}
// // 		next.ServeHTTP(w, r)
// // 	})
// // }

// // func main() {
// // 	http.HandleFunc("/api/getTkeyPubKey", getTkeyPubKeyHandler)
// // 	http.HandleFunc("/api/sign", signHandler)

// // 	enableCORS(mux)

// // 	fmt.Println("Server running on http://localhost:8080")
// // 	http.ListenAndServe(":8080", nil)
// // }

// // func getTkeyPubKeyHandler(w http.ResponseWriter, r *http.Request) {
// // 	pubKey, err := tkey.GetTkeyPubKey()
// // 	if err != nil {
// // 		http.Error(w, fmt.Sprintf("Failed to get public key: %v", err), http.StatusInternalServerError)
// // 		return
// // 	}

// // 	w.Header().Set("Content-Type", "application/json")
// // 	json.NewEncoder(w).Encode(map[string]string{"publicKey": string(pubKey)})
// // }

// // func signHandler(w http.ResponseWriter, r *http.Request) {
// // 	body, err := ioutil.ReadAll(r.Body)
// // 	if err != nil {
// // 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// // 		return
// // 	}

// // 	sig, err := tkey.Sign(body)
// // 	if err != nil {
// // 		http.Error(w, fmt.Sprintf("Failed to sign message: %v", err), http.StatusInternalServerError)
// // 		return
// // 	}

// // 	w.Header().Set("Content-Type", "application/json")
// // 	json.NewEncoder(w).Encode(map[string]string{"signature": string(sig)})
// // }

// package tkey

// import (
// 	"chalmers/tkey-group22/internal/tkey"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// )

// // CORS middleware
// func enableCORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 		if r.Method == "OPTIONS" {
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// func main() {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/api/getTkeyPubKey", getTkeyPubKeyHandler)
// 	mux.HandleFunc("/api/sign", signHandler)

// 	// Apply the CORS middleware
// 	handler := enableCORS(mux)

// 	fmt.Println("Server running on http://localhost:8080")
// 	http.ListenAndServe(":8080", handler)
// }

// func getTkeyPubKeyHandler(w http.ResponseWriter, r *http.Request) {
// 	pubKey, err := tkey.GetTkeyPubKey()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to get public key: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{"publicKey": string(pubKey)})
// }

// func signHandler(w http.ResponseWriter, r *http.Request) {
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	sig, err := tkey.Sign(body)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to sign message: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{"signature": string(sig)})
// }
