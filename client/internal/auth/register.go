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

const regurl = "http://localhost:8080/api/register"

func Register() error {

	username := getUsername()

	pubkey, err := tkey.GetTkeyPubKey()
	if err != nil {
		return err
	}

	sendRequest(pubkey, username)

	return nil
}

func getUsername() string {
	var username string
	fmt.Print("Please enter username: ")
	fmt.Scan(&username)
	return username
}

func sendRequest(pubkey ed25519.PublicKey, username string) {
	c := &http.Client{}

	data := RegisterRequest{Username: username, Pubkey: []byte(pubkey)}

	reqBody, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Post(regurl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Could not create user! Error: %s", res.Status)
		log.Fatal()
	} else {
		fmt.Printf("User '%s' has been successfully created!", username)
	}

}
