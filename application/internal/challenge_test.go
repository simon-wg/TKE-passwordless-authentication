package internal

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
	"time"
)

func TestVerifySignedResponse_ValidSignature(t *testing.T) {
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	pubKeyHex := hex.EncodeToString(pubKey)

	// Create an instance of ED25519ChallengeService
	service := NewED25519ChallengeService()

	// Generate a challenge
	challengeValue, err := service.GenerateChallenge(pubKeyHex)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes, _ := hex.DecodeString(challengeValue)
	signature := ed25519.Sign(privKey, challengeBytes)
	signatureHex := hex.EncodeToString(signature)

	// Verify the signed response
	valid, err := service.VerifySignature(pubKeyHex, signatureHex)
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

	// Create an instance of ED25519ChallengeService
	service := NewED25519ChallengeService()

	// Test with an invalid signature
	invalidSignature := []byte("invalidsignature")
	signatureHex := hex.EncodeToString(invalidSignature)
	valid, err := service.VerifySignature(pubKeyHex, signatureHex)
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

	// Create an instance of ED25519ChallengeService
	service := NewED25519ChallengeService()

	// Generate a challenge
	challengeValue, err := service.GenerateChallenge(pubKeyHex)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes, _ := hex.DecodeString(challengeValue)
	signature := ed25519.Sign(privKey, challengeBytes)
	signatureHex := hex.EncodeToString(signature)

	// Test with a non-existent challenge
	nonExistentPubKey := "nonexistentpubkey"
	valid, err := service.VerifySignature(nonExistentPubKey, signatureHex)
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
	if valid {
		t.Fatalf("Expected invalid signature for non-existent challenge, got valid")
	}
}

func TestVerifySignedResponse_ExpiredChallenge(t *testing.T) {
	service := NewED25519ChallengeService()

	// Set a new validDuration for the test
	service.validDuration = time.Duration(400) * time.Millisecond

	// // Restore the original validDuration after the test
	// defer func() {
	// 	validDuration = originalValidDuration
	// }()

	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	pubKeyHex := hex.EncodeToString(pubKey)

	// Create an instance of ED25519ChallengeService

	// Generate a challenge
	challengeValue, err := service.GenerateChallenge(pubKeyHex)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes, _ := hex.DecodeString(challengeValue)
	signature := ed25519.Sign(privKey, challengeBytes)
	signatureHex := hex.EncodeToString(signature)

	// Test with an expired challenge
	time.Sleep(service.validDuration + time.Duration(100)*time.Millisecond)
	valid, err := service.VerifySignature(pubKeyHex, signatureHex)
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
	if valid {
		t.Fatalf("Expected invalid signature for expired challenge, got valid")
	}
}
