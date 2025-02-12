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
	challenges      = make(map[string]*Challenge)
	validDuration   = 60 // challenges are valid for 60 seconds
	challengeLength = 64 // number of bytes in challenge
	challengesLock  sync.Mutex
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
		ExpiresAt: time.Now().Add(time.Duration(validDuration) * time.Second),
	}

	challenges[pubKey] = challenge

	return challenge.Value, nil
}

func cleanupExpiredChallenges() {
	for {
		time.Sleep(time.Duration(validDuration) * time.Second)
		challengesLock.Lock()
		for pubKey, challenge := range challenges {
			if time.Now().After(challenge.ExpiresAt) {
				delete(challenges, pubKey)
			}
		}
		challengesLock.Unlock()
	}
}
