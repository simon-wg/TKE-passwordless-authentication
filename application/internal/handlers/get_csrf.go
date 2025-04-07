package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func GetCSRF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
}
