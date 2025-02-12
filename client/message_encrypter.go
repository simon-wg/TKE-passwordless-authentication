package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func encrypt_message(message []byte, public_key *rsa.PublicKey, label []byte) []byte {
	hash := sha256.New()
	random := rand.Reader
	encrypted_message, err := rsa.EncryptOAEP(hash, random, public_key, message, label)
	if err != nil {
		panic(err)
	}

	return encrypted_message
}

func decrypt_message(ciphertext []byte, private_key *rsa.PrivateKey, label []byte) []byte {
	hash := sha256.New()
	random := rand.Reader
	decrypted_message, err := rsa.DecryptOAEP(hash, random, private_key, ciphertext, label)
	if err != nil {
		panic(err)
	}

	return decrypted_message
}
