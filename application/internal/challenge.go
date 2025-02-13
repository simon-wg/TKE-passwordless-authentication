package internal

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type Challenge struct {
	Value     string
	ExpiresAt time.Time
}

var (
	challengeLength  = 64                             // number of bytes in challenge
	validDuration    = time.Duration(1) * time.Minute // challenges are valid for 60 seconds
	cleanupInterval  = time.Duration(2) * time.Minute
	activeChallenges = make(map[string]*Challenge)
	challengesLock   sync.Mutex
)

func init() {
	go cleanupExpiredChallenges()
}

// Generate a new challenge for the given public key.
// It locks the challenges map to ensure thread safety, generates a random
// byte array, and encodes it to a hexadecimal string. The challenge is then
// stored in the challenges map with an expiration time.
//
// Parameters:
//   - pubKey: The public key for which the challenge is generated.
//
// Returns:
//   - A string representing the generated challenge value.
//
// Panics:
//   - If it fails to generate random bytes.
func GenerateChallenge(pubKey string) (string, error) {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	bytes := make([]byte, challengeLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.New("failed to generate random bytes")
	}

	challenge := &Challenge{
		Value:     hex.EncodeToString(bytes),
		ExpiresAt: time.Now().Add(validDuration),
	}

	activeChallenges[pubKey] = challenge

	return challenge.Value, nil
}

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
