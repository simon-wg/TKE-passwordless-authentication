package tests

import (
	"bytes"
	"chalmers/tkey-group22/application/internal"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test file could be changed to test all handlers.

const usersFilePath = "../data/users.csv"
const backupFilePath = "../data/users_backup.csv"
const loginURL = "/api/login"
const verifyURL = "/api/verify"
const mockUsername = "MockUser"

var mockPubKey ed25519.PublicKey
var mockPrivKey ed25519.PrivateKey

// backupUsersFile backs up the existing users.csv file to users_backup.csv.
//
// Parameters:
//   - t: The testing.T instance.
//
// Returns:
//   - None
func backupUsersFile(t *testing.T) {
	input, err := os.ReadFile(usersFilePath)
	if err != nil {
		t.Fatalf("Failed to read users file: %v", err)
	}

	err = os.WriteFile(backupFilePath, input, 0644)
	if err != nil {
		t.Fatalf("Failed to write backup file: %v", err)
	}
}

// restoreUsersFile restores the original users.csv file from the backup.
//
// Parameters:
//   - t: The testing.T instance.
//
// Returns:
//   - None
func restoreUsersFile(t *testing.T) {
	input, err := os.ReadFile(backupFilePath)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}

	err = os.WriteFile(usersFilePath, input, 0644)
	if err != nil {
		t.Fatalf("Failed to restore users file: %v", err)
	}

	err = os.Remove(backupFilePath)
	if err != nil {
		t.Fatalf("Failed to remove backup file: %v", err)
	}
}

// writeMockData writes mock data to the users.csv file.
//
// Parameters:
//   - t: The testing.T instance.
//   - data: A 2D slice of strings representing the mock data to be written.
//
// Returns:
//   - None
func writeMockData(t *testing.T, data [][]string) {
	file, err := os.Create(usersFilePath)
	if err != nil {
		t.Fatalf("Failed to create users file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			t.Fatalf("Failed to write record to users file: %v", err)
		}
	}
}

// TestMain sets up the test environment, runs the tests, and then restores the
// original environment.
//
// Parameters:
//   - m: The testing.M instance.
//
// Returns:
//   - None
//
// Dependencies:
//   - backupUsersFile
//   - writeMockData
//   - restoreUsersFile

func TestMain(m *testing.M) {
	backupUsersFile(nil)
	mockPubKey, mockPrivKey, _ = ed25519.GenerateKey(nil)
	mockData := [][]string{
		{"username", "publickey"},
		{"bob", "mocked_banana1234"},
		{"alice", "mocked_apple1234"},
		{mockUsername, string(mockPubKey)},
	}
	writeMockData(nil, mockData)

	// Run the tests including helper funcs
	code := m.Run()

	restoreUsersFile(nil)

	os.Exit(code)
}

// createRequest creates a new HTTP request and response recorder for testing.
//
// Parameters:
//   - t: The testing.T instance.
//   - method: The HTTP method.
//   - url: The URL for the request.
//   - body: A map representing the JSON body of the request.
//
// Returns:
//   - *httptest.ResponseRecorder: The response recorder for capturing the response.
//   - *http.Request: The created HTTP request.
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
	rr, req := createRequest(t, http.MethodPost, loginURL, map[string]string{"username": "bob"})
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
	rr, req := createRequest(t, http.MethodGet, loginURL, nil)
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

// Invalid req body. Expects fail.
func TestLoginHandler_InvalidRequestBody(t *testing.T) {
	rr, req := createRequest(t, http.MethodPost, loginURL, map[string]string{"invalid": "body"})
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// No username in req body. Expects fail.
func TestLoginHandler_UsernameNotProvided(t *testing.T) {
	rr, req := createRequest(t, http.MethodPost, loginURL, map[string]string{})
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// Bad user. Expects fail.
func TestLoginHandler_UserNotFound(t *testing.T) {
	rr, req := createRequest(t, http.MethodPost, loginURL, map[string]string{"username": "nonexistent"})
	internal.LoginHandler(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestVerifyHandler_InvalidRequestMethod(t *testing.T) {
	handler := http.HandlerFunc(internal.VerifyHandler)

	req, err := http.NewRequest(http.MethodGet, verifyURL, nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestVerifyHandler_InvalidRequestBody(t *testing.T) {
	handler := http.HandlerFunc(internal.VerifyHandler)

	req, err := http.NewRequest(http.MethodPost, verifyURL, bytes.NewBuffer([]byte("invalid body")))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestVerifyHandler_NonHexadecimalSignature(t *testing.T) {
	handler := http.HandlerFunc(internal.VerifyHandler)

	requestBody := map[string]string{
		"username":  mockUsername,
		"signature": "non-hex-signature",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, verifyURL, bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestVerifyHandler_NoActiveChallengeFound(t *testing.T) {
	handler := http.HandlerFunc(internal.VerifyHandler)

	// Generate a valid signature
	_, privKey, _ := ed25519.GenerateKey(nil)
	signature := ed25519.Sign(privKey, []byte("testChallenge"))

	requestBody := map[string]string{
		"username":  mockUsername,
		"signature": hex.EncodeToString(signature),
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, verifyURL, bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestVerifyHandler_InvalidSignature(t *testing.T) {
	handler := http.HandlerFunc(internal.VerifyHandler)

	internal.GenerateChallenge(mockUsername)

	// Generate random byte slice
	invalidSignBytes := make([]byte, 32)
	rand.Read(invalidSignBytes)

	requestBody := map[string]string{
		"username":  mockUsername,
		"signature": hex.EncodeToString(invalidSignBytes),
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, verifyURL, bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestVerifyHandler_VerificationSuccessful(t *testing.T) {
	handler := http.HandlerFunc(internal.VerifyHandler)

	challengeHex, _ := internal.GenerateChallenge(mockUsername)
	challengeBytes, _ := hex.DecodeString(challengeHex)
	signature := ed25519.Sign(mockPrivKey, challengeBytes)
	signatureHex := hex.EncodeToString(signature)

	requestBody := map[string]string{
		"username":  mockUsername,
		"signature": signatureHex,
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, verifyURL, bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
