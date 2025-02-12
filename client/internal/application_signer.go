package internal

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func GetSignatureAndVerify(user string) bool {
	pubkey := getPublicKey()

	sig, msg := getSignatureAndMessageFromEndpoint(user)

	return verifySig(pubkey, msg, sig)
}

func getPublicKey() []byte {
	baseUrl := "http://localhost:8080"
	endpoint := "/api/public"

	c := &http.Client{}

	resp, err := c.Get(baseUrl + endpoint)
	if err != nil {
		fmt.Println("Error sending request to get public key")
		return nil
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error in response from getting public key")
		return nil
	}

	var pubkey []byte

	err = json.NewDecoder(resp.Body).Decode(&pubkey)
	if err != nil {
		fmt.Println("Error with decoding public key")
		return nil
	}

	return pubkey
}

func verifySig(pubkey []byte, data []byte, sig []byte) bool {
	return ed25519.Verify(pubkey, data, sig)
}

func getSignatureAndMessageFromEndpoint(user string) ([]byte, []byte) {
	baseUrl := "http://localhost:8080"
	endpoint := "/api/login"

	// Get the signature and message from the endpoint

	httpClient := &http.Client{}

	resp, err := httpClient.PostForm(baseUrl+endpoint, url.Values{"user": {user}})
	if err != nil {
		fmt.Println("Error sending request to get signature and message")
		return nil, nil
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error in response from getting signature and message")
		return nil, nil
	}

	var sigMsg SignatureMessage
	err = json.NewDecoder(resp.Body).Decode(&sigMsg)

	if err != nil {
		fmt.Println("Error with decoding signature and message")
		return nil, nil
	}

	return sigMsg.Signature, sigMsg.Message
}

type SignatureMessage struct {
	Signature []byte
	Message   []byte
}
