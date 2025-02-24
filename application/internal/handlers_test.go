package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChallengeService struct {
	mock.Mock
}

func (m *MockChallengeService) HasActiveChallenge(publicKey string) bool {
	args := m.Called(publicKey)
	return args.Bool(0)
}

func (m *MockChallengeService) VerifySignature(publicKey, signature string) (bool, error) {
	args := m.Called(publicKey, signature)
	return args.Bool(0), args.Error(1)
}

func (m *MockChallengeService) GenerateChallenge(publicKey string) (string, error) {
	args := m.Called(publicKey)
	return args.String(0), args.Error(1)
}

func TestVerifyHandler_InvalidRequestMethod(t *testing.T) {
	mockChallengeService := new(MockChallengeService)
	handler := VerifyHandler(mockChallengeService)

	req, err := http.NewRequest(http.MethodGet, "/verify", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestVerifyHandler_InvalidRequestBody(t *testing.T) {
	mockChallengeService := new(MockChallengeService)
	handler := VerifyHandler(mockChallengeService)

	req, err := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer([]byte("invalid body")))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestVerifyHandler_NoActiveChallengeFound(t *testing.T) {
	mockChallengeService := new(MockChallengeService)
	handler := VerifyHandler(mockChallengeService)

	requestBody := map[string]string{
		"publicKey": "testPublicKey",
		"signature": "testSignature",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(body))
	assert.NoError(t, err)

	mockChallengeService.On("HasActiveChallenge", "testPublicKey").Return(false)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockChallengeService.AssertExpectations(t)
}

func TestVerifyHandler_InvalidSignature(t *testing.T) {
	mockChallengeService := new(MockChallengeService)
	handler := VerifyHandler(mockChallengeService)

	requestBody := map[string]string{
		"publicKey": "testPublicKey",
		"signature": "testSignature",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(body))
	assert.NoError(t, err)

	mockChallengeService.On("HasActiveChallenge", "testPublicKey").Return(true)
	mockChallengeService.On("VerifySignature", "testPublicKey", "testSignature").Return(false, nil)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockChallengeService.AssertExpectations(t)
}

func TestVerifyHandler_VerificationSuccessful(t *testing.T) {
	mockChallengeService := new(MockChallengeService)
	handler := VerifyHandler(mockChallengeService)

	requestBody := map[string]string{
		"publicKey": "testPublicKey",
		"signature": "testSignature",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(body))
	assert.NoError(t, err)

	mockChallengeService.On("HasActiveChallenge", "testPublicKey").Return(true)
	mockChallengeService.On("VerifySignature", "testPublicKey", "testSignature").Return(true, nil)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockChallengeService.AssertExpectations(t)
}
