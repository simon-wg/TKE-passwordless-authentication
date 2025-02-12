package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func register() error {
	pubkey, err := GetTkeyPubKey()
	if err != nil {
		return err
	}
	sendRequest(pubkey, "my_username")
	return nil
}

func getURL() string {
	url := "http://localhost:8080/api/register"
	return url
}

func convertToJSON(data map[string]string) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func createRequest(jsonData []byte, url string) *http.Request {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func sendRequest(pubkey []byte, username string) {

	url := getURL()

	pubkeyStr := string(pubkey[:])
	data := map[string]string{"username": username, "pubkey": pubkeyStr}

	jsonData := convertToJSON(data)

	req := createRequest(jsonData, url)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

}
