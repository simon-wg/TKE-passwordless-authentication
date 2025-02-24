package main

import (
	"chalmers/tkey-group22/internal/tkey"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/api/getTkeyPubKey", enableCors(http.HandlerFunc(getTkeyPubKeyHandler)))
	http.Handle("/api/sign", enableCors(http.HandlerFunc(signHandler)))

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

func getTkeyPubKeyHandler(w http.ResponseWriter, r *http.Request) {
	pubKey, err := tkey.GetTkeyPubKey()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get public key: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"publicKey": string(pubKey)})
}

func signHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	challengeHex, ok := requestBody["challenge"]
	if !ok {
		http.Error(w, "Missing challenge field", http.StatusBadRequest)
		return
	}

	challenge, err := hex.DecodeString(challengeHex)
	if err != nil {
		http.Error(w, "Invalid challenge format", http.StatusBadRequest)
		return
	}

	sig, err := tkey.Sign(challenge)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to sign message: %v", err), http.StatusInternalServerError)
		return
	}

	signatureHex := hex.EncodeToString(sig)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"signature": signatureHex})
}
