package internal

import (
<<<<<<< HEAD
=======
	"crypto"
>>>>>>> origin/dev/client
	"crypto/ed25519"
	"fmt"
	"log"
	"os"

	"github.com/tillitis/tkeyclient"
	"golang.org/x/crypto/ssh"
)

const progname = "tkey-device-signer"

var le = log.New(os.Stderr, "Error: ", 0)

<<<<<<< HEAD
func GetTkeyPubKey() []byte {
	devPath, err := tkeyclient.DetectSerialPort(false)
	if err != nil {
		fmt.Println("Error detecting serial port")
		return nil
=======
func GetTkeyPubKey() ([]byte, error) {
	signer, err := getSigner()

	if err != nil {
		return nil, err
	}

	pub, err := signer.tkSigner.GetPubkey()

	if err != nil {
		return nil, err
	}

	sshPub, _ := ssh.NewPublicKey(ed25519.PublicKey(pub))

	return ssh.MarshalAuthorizedKey(sshPub), nil
}

func Sign(msg []byte) ([]byte, error) {

	signer, err := getSigner()

	if err != nil {
		return nil, err
	}

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

func getSigner() (*Signer, error) {
	devPath, err := tkeyclient.DetectSerialPort(false)
	if err != nil {
		return nil, err
>>>>>>> origin/dev/client
	}

	serialSpeed := tkeyclient.SerialSpeed

	exit := func(code int) {
<<<<<<< HEAD
		fmt.Println("Error connecting to TKEY")
=======
>>>>>>> origin/dev/client
		os.Exit(0)
	}

	signer := NewSigner(devPath, serialSpeed, false, "", "", exit)

<<<<<<< HEAD
	pub, err := signer.tkSigner.GetPubkey()

	if err != nil {
		fmt.Println("Error getting Public Key")
		return nil
	}

	sshPub, _ := ssh.NewPublicKey(ed25519.PublicKey(pub))

	return ssh.MarshalAuthorizedKey(sshPub)
=======
	return signer, nil
>>>>>>> origin/dev/client
}
