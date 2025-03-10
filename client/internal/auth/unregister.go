package auth

import (
	"bytes"
	"chalmers/tkey-group22/client/internal/structs"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//
// Unregister firstly verifies the the user. If the user is verified, it calls sendUnregisterRequest to send an HTTP POST request to unregister the user.
//
// Parameters:
//   - username: The username of the user to be unregistered.
//	 - appurl: The URL of the application.
//
// Returns:
//   - nil: If the user is successfully unregistered.
//

func Unregister(appurl string, username string) error {

	challengeResponse, err := getChallenge(appurl, username)
	if err != nil {
		return err
	}

	signedChallenge, err := signChallenge(username, challengeResponse)
	if err != nil {
		return err
	}

	err = VerifyUser(appurl, username, signedChallenge)

	if err != nil {
		return fmt.Errorf("could not unregister user:  %s Verification failed", username)
	}

	sendUnregisterRequest(username)

	return nil
}

// sendUnregisterRequest sends an HTTP POST request to unregister a user.
//
// It constructs a JSON payload with the username and sends it to the unregistration endpoint.
// Depending on the HTTP response status code, it logs the appropriate message.
//
// Parameters:
//   - username: The username of the user to be unregistered.
//
// Returns:
//   - None
//

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
