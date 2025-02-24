package internal

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
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
func GenerateChallenge(pubkey []byte) (string, error) {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	bytes := make([]byte, challengeLength)
	rand.Read(bytes)

	challenge := &Challenge{
		Value:     hex.EncodeToString(bytes),
		ExpiresAt: time.Now().Add(ValidDuration),
	}

	activeChallenges[string(pubkey)] = challenge

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
func VerifySignature(pubkey []byte, signature []byte) (bool, error) {
	challenge, exists := activeChallenges[string(pubkey)]
	if !exists {
		return false, errors.New("no active challenge found for given key")
	}

	if time.Now().After(challenge.ExpiresAt) {
		return false, errors.New("challenge expired")
	}

	challengeBytes, err := hex.DecodeString(challenge.Value)
	if err != nil {
		return false, errors.New("invalid challenge format")
	}

	pubkey, _ = extractKeyBits(pubkey)
	if ed25519.Verify(ed25519.PublicKey(pubkey), challengeBytes, signature) {
		return true, nil
	}

	return false, errors.New("invalid signature")
}

func HasActiveChallenge(pubkey []byte) bool {
	challengesLock.Lock()
	defer challengesLock.Unlock()

	_, exists := activeChallenges[string(pubkey)]
	return exists
}

// extractKeyBits extracts the key bits from an SSH public key in byte slice format.
// The function expects the public key to be in the "ssh-ed25519" format and encoded in Base64.
// It returns the key bits as a byte slice or an error if the public key format is invalid or unsupported.
//
// Parameters:
//   - pubkey: A byte slice containing the SSH public key.
//
// Returns:
//   - A byte slice containing the extracted key bits.
//   - An error if the public key format is invalid, unsupported, or if there is an issue with Base64 decoding.
func extractKeyBits(pubkey []byte) ([]byte, error) {
	pubkeyString := string(pubkey)
	parts := strings.Split(pubkeyString, " ")
	if len(parts) < 2 {
		return nil, errors.New("invalid SSH public key format")
	}

	keyType := parts[0]
	if keyType != "ssh-ed25519" {
		return nil, errors.New("unsupported key type")
	}

	keyBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid Base64 encoding")
	}

	// The actual key bytes start after the key type and length prefix
	if len(keyBytes) < 32 {
		return nil, errors.New("invalid key length")
	}

	return keyBytes[len(keyBytes)-32:], nil
}
