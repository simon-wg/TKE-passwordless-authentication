package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// GetCSRF returns the a csrf token in the custom header X-CSRF-Token
// Access-Control-Expose-Headers allows the custom header to be read

func GetCSRF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")

}
