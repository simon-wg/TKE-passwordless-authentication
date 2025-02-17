package internal

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type ChallengeService interface {
	HasActiveChallenge(pubKey string) bool
	VerifySignature(pubKey string, signature string) (bool, error)
	GenerateChallenge(pubKey string) (string, error)
}

type Challenge struct {
	Value     string
	ExpiresAt time.Time
}

type ED25519ChallengeService struct {
	challengeLength  int           // Number of bytes in challenge
	validDuration    time.Duration // Time before the challenge is timed out
	cleanupInterval  time.Duration // Interval between removal of registered timed out challenges
	activeChallenges map[string]*Challenge
	challengesLock   sync.Mutex
}

func NewED25519ChallengeService() *ED25519ChallengeService {
	service := &ED25519ChallengeService{
		challengeLength:  128,
		validDuration:    time.Duration(1) * time.Minute,
		cleanupInterval:  time.Duration(2) * time.Minute,
		activeChallenges: make(map[string]*Challenge),
		challengesLock:   sync.Mutex{},
	}

	go service.cleanupExpiredChallenges()

	return service
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
func (s *ED25519ChallengeService) GenerateChallenge(pubKey string) (string, error) {
	s.challengesLock.Lock()
	defer s.challengesLock.Unlock()

	bytes := make([]byte, s.challengeLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.New("failed to generate random bytes")
	}

	challenge := &Challenge{
		Value:     hex.EncodeToString(bytes),
		ExpiresAt: time.Now().Add(s.validDuration),
	}

	s.activeChallenges[pubKey] = challenge

	return challenge.Value, nil
}

func (s *ED25519ChallengeService) cleanupExpiredChallenges() {
	for {
		time.Sleep(s.cleanupInterval)
		s.challengesLock.Lock()
		for pubKey, challenge := range s.activeChallenges {
			if time.Now().After(challenge.ExpiresAt) {
				delete(s.activeChallenges, pubKey)
			}
		}
		s.challengesLock.Unlock()
	}
}

// VerifySignature verifies the signed response for a given public key.
// It checks if there is an active challenge for the provided public key, if the challenge has not expired,
// and if the provided signature is valid for the challenge.
//
// Parameters:
//   - pubKey: The public key as a hexadecimal string.
//   - signature: The signature as a hexadecimal string.
//
// Returns:
//   - bool: True if the signature is valid, false otherwise.
//   - error: An error if the verification fails due to an invalid format, expired challenge, or no active challenge.
func (s *ED25519ChallengeService) VerifySignature(pubKey string, signature string) (bool, error) {
	challenge, exists := s.activeChallenges[pubKey]
	if !exists {
		return false, errors.New("no active challenge found for given key")
	}

	if time.Now().After(challenge.ExpiresAt) {
		return false, errors.New("challenge expired")
	}

	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return false, errors.New("invalid public key format")
	}

	challengeBytes, err := hex.DecodeString(challenge.Value)
	if err != nil {
		return false, errors.New("invalid challenge format")
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, errors.New("invalid signature format")
	}

	if ed25519.Verify(ed25519.PublicKey(pubKeyBytes), challengeBytes, signatureBytes) {
		return true, nil
	}

	return false, errors.New("invalid signature")
}

func (s *ED25519ChallengeService) HasActiveChallenge(pubKey string) bool {
	s.challengesLock.Lock()
	defer s.challengesLock.Unlock()

	_, exists := s.activeChallenges[pubKey]
	return exists
}
