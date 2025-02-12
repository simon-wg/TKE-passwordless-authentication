package internal

import (
	"fmt"
)

func GetChallengeAndSign(username string) ([]byte, error) {
	challenge, err := GetChallengeAndVerify(username)

	if err != nil {
		fmt.Println("Error getting challenge and verifying")
		return nil, err
	}

	// Sign the challenge
	sig, err := Sign(challenge.Message)
	if err != nil {
		fmt.Println("Error signing challenge")
		return nil, err
	}

	return sig, nil
}
