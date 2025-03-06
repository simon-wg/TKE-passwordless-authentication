package internal

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// Challenge represents a challenge that is generated for a user.
// It contains a random value and an expiration time.
type Challenge struct {
	Value     string
	ExpiresAt time.Time
}

var (
	ValidDuration    = time.Duration(1) * time.Minute // challenges are valid for 60 seconds
	challengeLength  = 128                            // number of bytes in challenge
	cleanupInterval  = time.Duration(2) * time.Minute
	activeChallenges = make(map[string]*Challenge)
	challengesLock   sync.Mutex
)

// Asynchronous function that runs in the background to clean up expired challenges.
// Asynchronous function that runs in the background to clean up expired challenges.
func init() {
	go cleanupExpiredChallenges()
}

// GenerateChallenge generates a new challenge for the given public key
// It creates a random byte sequence, encodes it to a hexadecimal string and stores it in the activeChallenges map with an expiration time
//
// Parameters:
//   - username: The public key for which the challenge is generated.
//
// Returns:
//   - A string representing the generated challenge.
//   - An error if the random byte generation fails.
func GenerateChallenge(username string) (string, error) {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	bytes := make([]byte, challengeLength)
	rand.Read(bytes)

	challenge := &Challenge{
		Value:     hex.EncodeToString(bytes),
		ExpiresAt: time.Now().Add(ValidDuration),
	}

	activeChallenges[username] = challenge

	return challenge.Value, nil
}

// cleanupExpiredChallenges periodically removes expired challenges from the activeChallenges map
// The function runs indefinitely, sleeping for a duration specified by cleanupInterval between each cleanup cycle
func cleanupExpiredChallenges() {
	for {
		time.Sleep(cleanupInterval)
		challengesLock.Lock()
		for pubKey, challenge := range activeChallenges {
			if time.Now().After(challenge.ExpiresAt) {
				delete(activeChallenges, pubKey)
			}
		}
		challengesLock.Unlock()
	}
}

// VerifySignature verifies the signed response for a given public key
//
// Parameters:
//   - username: The username as a string.
//   - signature: The signature as a byte slice.
//
// Returns:
//   - bool: True if the signature is valid, false otherwise.
//   - error: An error if the verification fails due to an invalid format, expired challenge, or no active challenge.
func VerifySignature(username string, signature []byte) (bool, error) {
	challenge, exists := activeChallenges[username]
	if !exists {
		return false, errors.New("no active challenge found for given user")
	}

	if time.Now().After(challenge.ExpiresAt) {
		return false, errors.New("challenge expired")
	}

	userData, err := UserRepo.GetUser(username)
	if err != nil {
		return false, err
	}

	for _, encodedPubKey := range userData.PublicKeys {
		pubKeyBytes, err := base64.StdEncoding.DecodeString(encodedPubKey)
		if err != nil {
			return false, err
		}
		edPubKey := ed25519.PublicKey(pubKeyBytes)
		if ed25519.Verify(edPubKey, []byte(challenge.Value), signature) {
			return true, nil
		}
	}

	return false, errors.New("invalid signature")
}

// HasActiveChallenge checks if there is an active challenge for the given user.
// It locks the challengesLock mutex to ensure thread safety while accessing the activeChallenges map.
//
// Parameters:
//   - username: The username to check for an active challenge.
//
// Returns:
//   - bool: True if there is an active challenge for the user, false otherwise.
func HasActiveChallenge(username string) bool {
	// Aquire the lock
	challengesLock.Lock()
	// Release the lock when the function returns
	defer challengesLock.Unlock()

	_, exists := activeChallenges[username]
	return exists
}
