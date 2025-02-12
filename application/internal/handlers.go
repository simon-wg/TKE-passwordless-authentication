package internal

// Handlers for register, login and verify

import (
	"crypto/rand"
    "encoding/hex"
    "sync"
    "time"
)

type Challenge struct {
	Value		string
	ExpiresAt	time.Time
}

var (
	challenges	    = make(map[string]*Challenge)
	validDuration   = 60 // challenges are valid for 60 seconds
	challengeLength = 64 // number of bytes in challenge
	challengesLock sync.Mutex
)

func init() {
	go cleanupExpiredChallenges()
}


func GenerateChallenge(pubKey string) string {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	bytes := make([]byte, 64)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("Failed to generate random bytes")
	}

	challenge := &Challenge {
		Value: 		hex.EncodeToString(bytes),
		ExpiresAt: 	ime.Now().Add(validDuration * time.Minute),
	}

	challenges[pubKey] = challenge

	return challenge.Value
}

func VerifyChallenge(pubKey string, signedMessage string) bool {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	challenge, exists := challenges[pubKey]
	if !exists || time.Now().After(challenge.ExpiresAt) {
		return false
	}

	// TODO: Implement actual verification of signedMessage with the challenge.Value

	return true
}

func cleanupExpiredChallenges() {
    for {
        time.Sleep(validDuration * time.seconds)
        challengesLock.Lock()
        for pubKey, challenge := range challenges {
            if time.Now().After(challenge.ExpiresAt) {
                delete(challenges, pubKey)
            }
        }
        challengesLock.Unlock()
    }
}
