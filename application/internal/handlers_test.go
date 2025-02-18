package internal

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyHandler_InvalidRequestMethod(t *testing.T) {
	handler := http.HandlerFunc(VerifyHandler)

	req, err := http.NewRequest(http.MethodGet, "/verify", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestVerifyHandler_InvalidRequestBody(t *testing.T) {
	handler := http.HandlerFunc(VerifyHandler)

	req, err := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer([]byte("invalid body")))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestVerifyHandler_NoActiveChallengeFound(t *testing.T) {
	handler := http.HandlerFunc(VerifyHandler)

	requestBody := map[string]string{
		"publicKey": "testPublicKey",
		"signature": "testSignature",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestVerifyHandler_InvalidSignature(t *testing.T) {
	handler := http.HandlerFunc(VerifyHandler)

	GenerateChallenge("testPublicKey")

	requestBody := map[string]string{
		"publicKey": "testPublicKey",
		"signature": "testSignature",
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestVerifyHandler_VerificationSuccessful(t *testing.T) {
	handler := http.HandlerFunc(VerifyHandler)

	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	pubKeyHex := hex.EncodeToString(pubKey)

	challengeHex, _ := GenerateChallenge(pubKeyHex)
	challengeBytes, _ := hex.DecodeString(challengeHex)
	signature := ed25519.Sign(privKey, challengeBytes)
	signatureHex := hex.EncodeToString(signature)

	requestBody := map[string]string{
		"publicKey": pubKeyHex,
		"signature": signatureHex,
	}
	body, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
