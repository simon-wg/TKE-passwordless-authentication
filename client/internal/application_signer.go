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
		fmt.Println("Error getting public key")
		return nil, err
	}

	sigMsg, err := fetchMessageAndSignature(user)
	if err != nil {
		fmt.Println("Error fetching message and signature")
		return nil, err
	}

	if !verifySig(pubkey, *sigMsg) {
		fmt.Println("Signature verification failed")
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
		fmt.Println("Error sending request to get public key")
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error in response from getting public key")
		return nil, fmt.Errorf("error in response")
	}

	var pubkey []byte

	err = json.NewDecoder(resp.Body).Decode(&pubkey)
	if err != nil {
		fmt.Println("Error with decoding public key")
		return nil, err
	}

	return pubkey, nil
}

func verifySig(pubkey ed25519.PublicKey, sigMsg SignatureMessage) bool {
	if len(pubkey) != ed25519.PublicKeySize {
		fmt.Println("Invalid public key")
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
		fmt.Println("Error sending request to get signature and message")
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error in response from getting signature and message")
		return nil, fmt.Errorf("error in response")
	}

	var sigMsg SignatureMessage
	err = json.NewDecoder(resp.Body).Decode(&sigMsg)

	if err != nil {
		fmt.Println("Error with decoding signature and message")
		return nil, fmt.Errorf("error decoding")
	}

	return &sigMsg, nil
}

type SignatureMessage struct {
	Signature []byte
	Message   []byte
}
