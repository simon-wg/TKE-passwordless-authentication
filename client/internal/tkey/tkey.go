package tkey

import (
	"bufio"
	"crypto"
	"crypto/ed25519"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tillitis/tkeyclient"
)

const progname = "tkey-device-signer"

var le = log.New(os.Stderr, "Error: ", 0)
var existingSigner *Signer

// GetTkeyPubKey retrieves the public key from the TKey signer.
// It first obtains a signer instance and attempts to connect to it.
// If the connection is successful, it fetches the public key from the signer,
// prints the authorized key, and returns the public key.
// If any error occurs during these steps, it returns the error.
//
// Returns:
//   - ed25519.PublicKey: The public key retrieved from the TKey signer.
//   - error: An error if any step fails, otherwise nil.
func GetTkeyPubKey() (ed25519.PublicKey, error) {
	signer, err := getSigner()

	if err != nil {
		return nil, err
	}

	if !signer.connect() {
		le.Printf("Connect failed")
		return nil, fmt.Errorf("connect failed")
	}

	defer signer.disconnect()

	pub, err := signer.tkSigner.GetPubkey()

	if err != nil {
		return nil, err
	}

	pubkey := ed25519.PublicKey(pub)

	// TODO: Remove this later. This is just for testing purposes.
	signer.printAuthorizedKey()

	return pubkey, nil
}

// Sign signs the given message using a signer obtained from getSigner.
// It returns the signature or an error if the signing process fails.
//
// Parameters:
//   - msg: The message to be signed.
//
// Returns:
//   - []byte: The generated signature.
//   - error: An error if the signing process fails.
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

// getSigner initializes and returns a Signer instance. If an existing signer is already
// connected and is the desired application, it returns the existing signer. Otherwise,
// it detects the serial port, prompts the user to enter a User Supplied Secret (USS)
// manually or provide a USS file, and creates a new Signer instance.
//
// Returns:
//   - *Signer: A pointer to the initialized Signer instance.
//   - error: An error if the signer could not be initialized or the serial port could not be detected.
func getSigner() (*Signer, error) {
	if existingSigner != nil && existingSigner.connect() && existingSigner.isWantedApp() {
		// The signer app is already loaded, return the existing signer
		return existingSigner, nil
	}

	devPath, err := tkeyclient.DetectSerialPort(false)
	if err != nil {
		return nil, err
	}

	serialSpeed := tkeyclient.SerialSpeed

	exit := func(code int) {
		os.Exit(0)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to enter a manual User Supplied Secret (USS) or provide a USS file? (m/f/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	enterUSS := false
	fileUSS := ""

	if response == "m" {
		enterUSS = true
	} else if response == "f" {
		fmt.Print("Please provide the path to the USS file: ")
		fileUSS, _ = reader.ReadString('\n')
		fileUSS = strings.TrimSpace(fileUSS)
	}

	signer := NewSigner(devPath, serialSpeed, enterUSS, fileUSS, "", exit)
	existingSigner = signer

	return signer, nil
}
