package internal

import (
	"fmt"
	"net/http"
	"net/url"
)

func getChallengeAndSign(username string) ([]byte, error) {
	challenge, err := GetChallengeAndVerify(username)

	if err != nil {
		return nil, err
	}

	// Sign the challenge
	sig, err := Sign(challenge.Message)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func Login(username string) error {
	c := &http.Client{}

	sig, err := getChallengeAndSign(username)
	if err != nil {
		return err
	}

	baseUrl := "http://localhost:8080"
	endpoint := "/api/verify"

	resp, err := c.PostForm(baseUrl+endpoint, url.Values{"username": {username}, "signature": {string(sig)}})

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error in response when sending signature")
	}

	return nil
}
