package util

import (
	"crypto/ed25519"
	"crypto/rand"
	"log"
)

type KeyInfo struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func GenerateNewEd25519Keys() (*KeyInfo, error) {
	// Generate a new private/public key pair using Ed25519
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Error generating Ed25519 key pair: %v", err)
		return nil, err
	}

	return &KeyInfo{
		PrivateKey: privKey,
		PublicKey:  pubKey,
	}, nil
}
