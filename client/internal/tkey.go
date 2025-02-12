package internal

import (
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

	pub, err := signer.tkSigner.GetPubkey()

	if err != nil {
		fmt.Println("Error getting Public Key")
		return nil
	}

	sshPub, _ := ssh.NewPublicKey(ed25519.PublicKey(pub))

	return ssh.MarshalAuthorizedKey(sshPub)
}
