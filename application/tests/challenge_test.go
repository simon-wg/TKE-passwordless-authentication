package tests

import (
	"chalmers/tkey-group22/application/internal"
	"crypto/ed25519"
	"encoding/hex"
	"testing"
	"time"
)

func TestVerifySignedResponse_ValidSignature(t *testing.T) {
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	pubKeyHex := hex.EncodeToString(pubKey)

	// Generate a challenge
	challengeValue, err := internal.GenerateChallenge(pubKeyHex)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes, _ := hex.DecodeString(challengeValue)
	signature := ed25519.Sign(privKey, challengeBytes)
	signatureHex := hex.EncodeToString(signature)

	// Verify the signed response
	valid, err := internal.VerifySignature(pubKeyHex, signatureHex)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !valid {
		t.Fatalf("Expected valid signature, got invalid")
	}
}

func TestVerifySignedResponse_InvalidSignature(t *testing.T) {
	pubKey, _, _ := ed25519.GenerateKey(nil)
	pubKeyHex := hex.EncodeToString(pubKey)

	invalidSignature := []byte("invalidsignature")
	signatureHex := hex.EncodeToString(invalidSignature)
	valid, err := internal.VerifySignature(pubKeyHex, signatureHex)
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
	if valid {
		t.Fatalf("Expected invalid signature, got valid")
	}
}

func TestVerifySignedResponse_NonExistentChallenge(t *testing.T) {
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	pubKeyHex := hex.EncodeToString(pubKey)

	// Generate a challenge
	challengeValue, err := internal.GenerateChallenge(pubKeyHex)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes, _ := hex.DecodeString(challengeValue)
	signature := ed25519.Sign(privKey, challengeBytes)
	signatureHex := hex.EncodeToString(signature)

	// Test with a non-existent challenge
	nonExistentPubKey := "nonexistentpubkey"
	valid, err := internal.VerifySignature(nonExistentPubKey, signatureHex)
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
	if valid {
		t.Fatalf("Expected invalid signature for non-existent challenge, got valid")
	}
}

func TestVerifySignedResponse_ExpiredChallenge(t *testing.T) {
	// Set a new validDuration for the test
	originalValidDuration := internal.ValidDuration
	internal.ValidDuration = time.Duration(400) * time.Millisecond

	// // Restore the original validDuration after the test
	defer func() {
		internal.ValidDuration = originalValidDuration
	}()

	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	pubKeyHex := hex.EncodeToString(pubKey)

	// Generate a challenge
	challengeValue, err := internal.GenerateChallenge(pubKeyHex)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes, _ := hex.DecodeString(challengeValue)
	signature := ed25519.Sign(privKey, challengeBytes)
	signatureHex := hex.EncodeToString(signature)

	// Test with an expired challenge
	time.Sleep(internal.ValidDuration + time.Duration(100)*time.Millisecond)
	valid, err := internal.VerifySignature(pubKeyHex, signatureHex)
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
	if valid {
		t.Fatalf("Expected invalid signature for expired challenge, got valid")
	}
}
