package internal

import (
	"crypto"
	"crypto/ed25519"
	"fmt"
	"log"
	"os"

	"github.com/tillitis/tkeyclient"
	"golang.org/x/crypto/ssh"
)

const progname = "tkey-device-signer"

var le = log.New(os.Stderr, "Error: ", 0)

func GetTkeyPubKey() []byte {
	signer := getSigner()

	pub, err := signer.tkSigner.GetPubkey()

	if err != nil {
		fmt.Println("Error getting Public Key")
		return nil
	}

	sshPub, _ := ssh.NewPublicKey(ed25519.PublicKey(pub))

	return ssh.MarshalAuthorizedKey(sshPub)
}

func Sign(msg []byte) ([]byte, error) {

	signer := getSigner()

	if !signer.connect() {
		le.Printf("Connect failed")
		return nil, fmt.Errorf("connect failed")
	}

	defer signer.disconnect()

	sig, err := signer.Sign(nil, msg, crypto.Hash(0))
	if err != nil {
		le.Printf("Sign failed: %s\n", err)
		return nil, err
	}

	return sig, nil
}

func getSigner() *Signer {
	devPath, err := tkeyclient.DetectSerialPort(false)
	if err != nil {
		fmt.Println("Error detecting serial port")
		return nil
	}

	serialSpeed := tkeyclient.SerialSpeed

	exit := func(code int) {
		fmt.Println("Error connecting to TKEY")
		os.Exit(0)
	}

	signer := NewSigner(devPath, serialSpeed, false, "", "", exit)

	return signer
}
