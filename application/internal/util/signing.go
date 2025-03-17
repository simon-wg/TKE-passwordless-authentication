package util

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func loadPrivateKeyFromEnv() (*rsa.PrivateKey, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Get the private key string
	encodedKey := os.Getenv("PRIVATE_KEY")
	if encodedKey == "" {
		return nil, fmt.Errorf("PRIVATE_KEY not set in .env")
	}

	// Decode from base64
	privKeyBytes, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 private key: %v", err)
	}

	// Parse PEM block
	block, _ := pem.Decode(privKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to parse PEM block containing private key")
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %v", err)
	}

	return privateKey, nil
}

func SignMessage(message []byte) ([]byte, error) {

	privateKey, err := loadPrivateKeyFromEnv()
	if err != nil {
		return nil, fmt.Errorf("Failed to load key")
	}
	// Hash the message
	hashed := sha256.Sum256(message)

	// Sign the hashed message
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	fmt.Println(signature)
	return signature, nil
}
