package tests

import (
	"chalmers/tkey-group22/application/internal"
	"chalmers/tkey-group22/application/internal/util"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validPubKey ed25519.PublicKey
var validPrivKey ed25519.PrivateKey
var testUsername = "TestUser"
var mockRepo = UserRepoMock{}

func init() {
	validPubKey, validPrivKey, _ = ed25519.GenerateKey(nil)
}

// Mock required UserRepo functions

type UserRepoMock struct{}

func (u UserRepoMock) GetUser(username string) (*util.User, error) {
	if username == testUsername {
		return &util.User{
			ID:         primitive.NewObjectID(),
			Username:   testUsername,
			PublicKeys: []util.PublicKey{util.PublicKey{Label: "main", Key: base64.StdEncoding.EncodeToString(validPubKey)}},
		}, nil
	}
	return nil, errors.New("No user found")
}

func (u UserRepoMock) AddPublicKey(userName string, newPubKey ed25519.PublicKey, label string) (*mongo.UpdateResult, error) {
	panic("unimplemented")
}

func (u UserRepoMock) CreateUser(userName string, pubkey ed25519.PublicKey, label string) (*mongo.InsertOneResult, error) {
	panic("unimplemented")
}

func (u UserRepoMock) DeleteUser(userName string) (*mongo.DeleteResult, error) {
	panic("unimplemented")
}

func (u UserRepoMock) GetPublicKeyLabels(userName string) ([]string, error) {
	panic("unimplemented")
}

func (u UserRepoMock) RemovePublicKey(userName string, label string) (*mongo.UpdateResult, error) {
	panic("unimplemented")
}

func (u UserRepoMock) UpdateUser(userName string, updatedUser util.User) (*mongo.UpdateResult, error) {
	panic("unimplemented")
}

func TestVerifySignedResponse_ValidSignature(t *testing.T) {
	username := "TestUser"

	// Generate a challenge
	challengeValue, err := internal.GenerateChallenge(username)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes := []byte(challengeValue)
	signature := ed25519.Sign(validPrivKey, challengeBytes)

	// Verify the signed response
	valid, err := internal.VerifySignature(username, signature, mockRepo)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !valid {
		t.Fatalf("Expected valid signature, got invalid")
	}
}

func TestVerifySignedResponse_InvalidSignature(t *testing.T) {

	invalidSignature := []byte("invalidsignature")
	valid, err := internal.VerifySignature(testUsername, invalidSignature, mockRepo)
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
	if valid {
		t.Fatalf("Expected invalid signature, got valid")
	}
}

func TestVerifySignedResponse_NonExistentChallenge(t *testing.T) {

	// Generate a challenge
	challengeValue, err := internal.GenerateChallenge(testUsername)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes := []byte(challengeValue)
	signature := ed25519.Sign(validPrivKey, challengeBytes)

	// Test with a non-existent challenge
	valid, err := internal.VerifySignature("nonExistingUser", signature, mockRepo)
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
	internal.ValidDuration = time.Duration(200) * time.Millisecond

	// // Restore the original validDuration after the test
	defer func() {
		internal.ValidDuration = originalValidDuration
	}()

	fmt.Println(internal.ValidDuration)
	// Generate a challenge
	challengeValue, err := internal.GenerateChallenge(testUsername)
	if err != nil {
		t.Fatalf("Failed to generate challenge: %v", err)
	}

	// Sign the challenge
	challengeBytes := []byte(challengeValue)
	signature := ed25519.Sign(validPrivKey, challengeBytes)

	// Sleep until challenge expires
	time.Sleep(internal.ValidDuration + time.Duration(100)*time.Millisecond)
	valid, err := internal.VerifySignature(testUsername, signature, mockRepo)
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
	if valid {
		t.Fatalf("Expected invalid signature for expired challenge, got valid")
	}
}
