package tkey

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/api/getTkeyPubKey", getTkeyPubKeyHandler)
	http.HandleFunc("/api/sign", signHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func getTkeyPubKeyHandler(w http.ResponseWriter, r *http.Request) {
	pubKey, err := GetTkeyPubKey()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get public key: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"publicKey": string(pubKey)})
}

func signHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sig, err := Sign(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to sign message: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"signature": string(sig)})
}
