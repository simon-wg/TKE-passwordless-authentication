package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tillitis/tkeyclient"
)

const progname = "tkey-device-signer"

var le = log.New(os.Stderr, "", 0)

func main() {
	devPath := "/dev/ttyACM0"

	// tk := tkeyclient.New()
	// err := tk.Connect(devPath)
	// if err != nil {
	// 	panic("Failed to connect to TKEY")
	// }
	// fmt.Printf("Connected to TKEY\n")

	// nameVer, err := tk.GetNameVersion()
	// if err != nil {
	// 	panic("Failed to get Name and Version")
	// }
	// fmt.Printf("Firmware name: %s\n", nameVer.Name0)
	// fmt.Printf("Name1: %s\n", nameVer.Name1)
	// fmt.Printf("Version: %d\n", nameVer.Version)

	exit := func(code int) {
		os.Exit(code)
	}

	signer := NewSigner(devPath, tkeyclient.SerialSpeed, false, "", "", exit)
	if !signer.connect() {
		le.Printf("Connect failed")
		return
	}

	defer signer.disconnect()

	fmt.Printf("Connected to TKEY\n")
	signer.printAuthorizedKey()

	sign, err := signer.Sign(os.Stdin, []byte("test"), nil)
	if err != nil {
		le.Printf("Sign failed: %s\n", err)
		return
	}
	fmt.Printf("Signature: %s\n", sign)
}
