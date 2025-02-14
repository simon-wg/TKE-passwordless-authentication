package tests

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"chalmers/tkey-group22/application/internal"
)

const usersFilePath = "../data/users.csv"
const backupFilePath = "../data/users_backup.csv"

// Helper function to backup the existing users.csv file
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

// Helper function to restore the original users.csv file from the backup
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

// Helper function to write mock data to users.csv file
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

func TestMain(m *testing.M) {
	// Backup the existing users.csv file
	backupUsersFile(nil)

	// Write mock data to users.csv file
	mockData := [][]string{
		{"username", "publickey"},
		{"bob", "mocked_banana1234"},
		{"alice", "mocked_apple1234"},
	}
	writeMockData(nil, mockData)

	// Run the tests
	code := m.Run()

	// Restore the original users.csv file
	restoreUsersFile(nil)

	os.Exit(code)
}

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
