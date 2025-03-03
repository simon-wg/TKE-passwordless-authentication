package auth

import (
	"bytes"
	. "chalmers/tkey-group22/internal/structs"
	"chalmers/tkey-group22/internal/tkey"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Sends a request to the server to login a user
// Returns an error if the request fails
func Login(appurl string, username string) error {
	c := &http.Client{}

	// Fetches the generated challenge from the server
	challengeResponse, err := getChallenge(appurl, username)
	if err != nil {
		return err
	}

	// TODO: Implement signature verification
	// if !verifySignature(challengeResponse) {
	// 	return fmt.Errorf("signature verification failed")
	// }

	// Signs the challenge
	signedChallenge, err := signChallenge(username, challengeResponse)
	if err != nil {
		return err
	}

	// TODO: Make more customizable
	endpoint := "/api/verify"

	body, err := json.Marshal(signedChallenge)
	if err != nil {
		return err
	}

	// Sends the signed challenge to the server in the format of a VerifyRequest
	resp, err := c.Post(appurl+endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("User '%s' has been successfully logged in!\n", username)
	case http.StatusUnauthorized:
		return fmt.Errorf("invalid signature")
	case http.StatusNotFound:
		return fmt.Errorf("no active challenge found for the user")
	case http.StatusInternalServerError:
		return fmt.Errorf("unable to read user data")
	default:
		return fmt.Errorf("unexpected error: %s", resp.Status)
	}

	return nil
}

// An internal function that signs the challenge using the tkey
func signChallenge(username string, challenge *LoginResponse) (*VerifyRequest, error) {
	fmt.Printf("Touch the TKey to continue...\n")
	sig, err := tkey.Sign([]byte(challenge.Challenge))
	if err != nil {
		return nil, err
	}

	return &VerifyRequest{
		Username:  username,
		Signature: sig,
	}, nil
}

// An internal function that fetches the challenge from the server
func getChallenge(appurl string, user string) (*LoginResponse, error) {
	endpoint := "/api/login"

	c := &http.Client{}

	body, err := json.Marshal(LoginRequest{Username: user})
	if err != nil {
		return nil, err
	}

	// Get the signature and message from the endpoint
	resp, err := c.Post(appurl+endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// Continue processing
	case http.StatusNotFound:
		return nil, fmt.Errorf("user not found")
	case http.StatusBadRequest:
		return nil, fmt.Errorf("invalid request body or missing username")
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("unable to read user data")
	default:
		return nil, fmt.Errorf("unexpected error: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body")
	}

	// Decode the response body into a LoginResponse struct
	var res *LoginResponse
	err = json.Unmarshal(respBody, &res)

	if err != nil {
		return nil, fmt.Errorf("error decoding challenge response")
	}

	return res, nil
}

// TODO: Uncomment once server verification is implemented
// func verifySignature(res *LoginResponse) bool {
// 	pubkey, err := getPublicKey()
// 	if err != nil {
// 		return false
// 	}
// 	return ed25519.Verify(pubkey, []byte(res.Challenge), []byte(res.Signature))
// }

// func getPublicKey() ([]byte, error) {
// 	baseUrl := "http://localhost:8080"
// 	endpoint := "/api/public"

// 	c := &http.Client{}

// 	resp, err := c.Get(baseUrl + endpoint)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("error in response when getting public key")
// 	}

// 	var pubkey []byte

// 	err = json.NewDecoder(resp.Body).Decode(&pubkey)

// 	return pubkey, err
// }
