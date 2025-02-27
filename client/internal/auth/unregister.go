package auth

import (
	"bytes"
	"chalmers/tkey-group22/internal/structs"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Unregister() error {

	username := GetUsername()
	err := VerifyUser(username)

	if err != nil {
		return fmt.Errorf("user verification failed for %s", username)
	}

	if err != nil {
		return fmt.Errorf("could not retrieve public key from tkey")

	}
	sendUnregisterRequest(username)

	return nil
}

func sendUnregisterRequest(username string) {
	const regurl = "http://localhost:8080/api/unregister"
	c := &http.Client{}

	data := structs.UnregisterRequest{Username: username}
	reqBody, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	res, err := c.Post(regurl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		fmt.Printf("User '%s' has successfully been unregistered!\n", username)
	case http.StatusConflict:
		fmt.Printf("User '%s' does not exists!\n", username)
	case http.StatusBadRequest:
		fmt.Printf("Invalid request body for user '%s'\n", username)
	case http.StatusInternalServerError:
		fmt.Printf("Unable to unregister user '%s'\n", username)
	default:
		fmt.Printf("Unexpected error: %s\n", res.Status)
	}

}
