package auth

import (
	"bytes"
	. "chalmers/tkey-group22/internal/structs"
	"chalmers/tkey-group22/internal/tkey"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// const regurl = "http://localhost:8080/api/register"

func Register(appurl string, username string) error {

	// username := GetUsername()

	pubkey, err := tkey.GetTkeyPubKey()
	if err != nil {
		return err
	}

	err = sendRequest(appurl, pubkey, username)
	if err != nil {
		return err
	}

	return nil
}

func sendRequest(appurl string, pubkey ed25519.PublicKey, username string) error {
	c := &http.Client{}

	data := RegisterRequest{Username: username, Pubkey: []byte(pubkey)}

	reqBody, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Post(appurl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		fmt.Printf("User '%s' has been successfully created!\n", username)
		return nil
	case http.StatusConflict:
		return fmt.Errorf("user '%s' already exists", username)
	case http.StatusBadRequest:
		return fmt.Errorf("invalid request body for user '%s'", username)
	case http.StatusInternalServerError:
		return fmt.Errorf("unable to save user data for user '%s'", username)
	default:
		return fmt.Errorf("unexpected error: %s", res.Status)
	}

}
