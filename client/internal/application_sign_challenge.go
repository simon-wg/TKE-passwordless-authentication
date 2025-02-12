package internal

import (
	"crypto"
	"fmt"
	"os"

	"github.com/tillitis/tkeyclient"
)

func GetChallengeAndSign(username string) []byte {
	challenge, verification := GetChallengeAndVerify(username)

	if !verification {
		fmt.Println("Verification failed")
		return nil
	}

	// Sign the challenge
	sig := signChallenge(challenge)

	if sig == nil {
		fmt.Println("Failed to sign challenge")
		return nil
	}

	return sig
}

func signChallenge(challenge []byte) []byte {
	devPath := getSerialPort()
	serialSpeed := tkeyclient.SerialSpeed

	exit := func(code int) {
		os.Exit(code)
	}

	signer := NewSigner(devPath, serialSpeed, false, "", "", exit)

	if !signer.connect() {
		le.Printf("Connect failed")
		return nil
	}

	defer signer.disconnect()

	sig, err := signer.Sign(nil, challenge, crypto.Hash(0))
	if err != nil {
		le.Printf("Sign failed: %s\n", err)
		return nil
	}

	return sig
}

func getSerialPort() string {
	devPath, err := tkeyclient.DetectSerialPort(false)
	if err != nil {
		le.Printf("Failed to detect serial port: %s\n", err)
		return ""
	}
	return devPath
}
