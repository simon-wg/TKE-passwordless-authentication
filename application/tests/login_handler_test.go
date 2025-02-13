package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chalmers/tkey-group22/application/internal"
)

// Helper function to create a new HTTP request and recorder
func createRequest(t *testing.T, method, url string, body map[string]string) (*httptest.ResponseRecorder, *http.Request) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	return rr, req
}

// Valid input. Expects success.
func TestLoginHandler_Success(t *testing.T) {
	rr, req := createRequest(t, http.MethodGet, "/api/login", map[string]string{"username": "bob"})
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"challenge":"`
	if !bytes.HasPrefix(rr.Body.Bytes(), []byte(expected)) {
		t.Errorf("handler returned unexpected body: got %v want prefix %v", rr.Body.String(), expected)
	}
}

// Invalid method. Expects fail.
func TestLoginHandler_InvalidMethod(t *testing.T) {
	rr, req := createRequest(t, http.MethodPost, "/api/login", nil)
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

// Invalid req body. Expects fail.
func TestLoginHandler_InvalidRequestBody(t *testing.T) {
	rr, req := createRequest(t, http.MethodGet, "/api/login", map[string]string{"invalid": "body"})
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// No username in req body. Expects fail.
func TestLoginHandler_UsernameNotProvided(t *testing.T) {
	rr, req := createRequest(t, http.MethodGet, "/api/login", map[string]string{})
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// Bad user. Expects fail.
func TestLoginHandler_UserNotFound(t *testing.T) {
	rr, req := createRequest(t, http.MethodGet, "/api/login", map[string]string{"username": "nonexistent"})
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
