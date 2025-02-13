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

// GenerateChallenge generates a new challenge for the given public key.
// It creates a random byte sequence, encodes it to a hexadecimal string,
// and stores it in the activeChallenges map with an expiration time.
//
// Parameters:
//   - pubKey: The public key for which the challenge is generated.
//
// Returns:
//   - A string representing the generated challenge.
//   - An error if the random byte generation fails.
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
