package internal

import (
	"chalmers/tkey-group22/application/internal/util"
	"crypto/ed25519"
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
	ValidDuration    = time.Duration(1) * time.Minute // challenges are valid for 60 seconds
	challengeLength  = 128                            // number of bytes in challenge
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
func GenerateChallenge(user string) (string, error) {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	bytes := make([]byte, challengeLength)
	rand.Read(bytes)

	challenge := &Challenge{
		Value:     hex.EncodeToString(bytes),
		ExpiresAt: time.Now().Add(ValidDuration),
	}

	activeChallenges[user] = challenge

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

// VerifySignature verifies the signed response for a given public key.
// It checks if there is an active challenge for the provided public key, if the challenge has not expired,
// and if the provided signature is valid for the challenge.
//
// Parameters:
//   - username: The username as a string.
//   - signature: The signature as a hexadecimal string.
//
// Returns:
//   - bool: True if the signature is valid, false otherwise.
//   - error: An error if the verification fails due to an invalid format, expired challenge, or no active challenge.
func VerifySignature(user string, signature []byte) (bool, error) {
	challenge, exists := activeChallenges[user]
	if !exists {
		return false, errors.New("no active challenge found for given key")
	}

	if time.Now().After(challenge.ExpiresAt) {
		return false, errors.New("challenge expired")
	}

	userData, err := util.Read(UsersFile)
	if err != nil {
		return false, errors.New("unable to read user data")
	}

	pubkeyString := userData[user]

	edPubkey := ed25519.PublicKey([]byte(pubkeyString))

	if ed25519.Verify(edPubkey, []byte(challenge.Value), signature) {
		return true, nil
	}

	return false, errors.New("invalid signature")
}

func HasActiveChallenge(user string) bool {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	_, exists := activeChallenges[user]
	return exists
}
