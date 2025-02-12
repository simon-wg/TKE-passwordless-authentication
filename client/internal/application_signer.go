package internal

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func GetChallengeAndVerify(user string) (*SignatureMessage, error) {
	pubkey, err := getPublicKey()
	if err != nil {
		return nil, err
	}

	sigMsg, err := fetchMessageAndSignature(user)
	if err != nil {
		return nil, err
	}

	if !verifySig(pubkey, *sigMsg) {
		return nil, fmt.Errorf("signature verification failed")
	}

	return sigMsg, nil
}

func getPublicKey() ([]byte, error) {
	baseUrl := "http://localhost:8080"
	endpoint := "/api/public"

	c := &http.Client{}

	resp, err := c.Get(baseUrl + endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error in response when getting public key")
	}

	var pubkey []byte

	err = json.NewDecoder(resp.Body).Decode(&pubkey)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

func verifySig(pubkey ed25519.PublicKey, sigMsg SignatureMessage) bool {
	if len(pubkey) != ed25519.PublicKeySize {
		return false
	}
	return ed25519.Verify(pubkey, sigMsg.Message, sigMsg.Signature)
}

func fetchMessageAndSignature(user string) (*SignatureMessage, error) {
	baseUrl := "http://localhost:8080"
	endpoint := "/api/login"

	// Get the signature and message from the endpoint

	c := &http.Client{}

	resp, err := c.PostForm(baseUrl+endpoint, url.Values{"user": {user}})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error in response when requesting challenge")
	}

	var sigMsg SignatureMessage
	err = json.NewDecoder(resp.Body).Decode(&sigMsg)

	if err != nil {
		return nil, fmt.Errorf("error decoding challenge response")
	}

	return &sigMsg, nil
}

type SignatureMessage struct {
	Signature []byte
	Message   []byte
}
