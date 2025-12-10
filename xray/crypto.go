package xray

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRealityX25519PrivateKey() (string, error) {
	curve := ecdh.X25519()

	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(priv.Bytes()), nil
}

func DeriveRealityX25519PublicKey(privateKey string) (string, error) {
	curve := ecdh.X25519()

	privBytes, err := base64.RawURLEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key encoding: %w", err)
	}
	if len(privBytes) != 32 {
		return "", fmt.Errorf("private key must be 32 bytes (got %d)", len(privBytes))
	}

	priv, err := curve.NewPrivateKey(privBytes)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	pub := priv.PublicKey()

	return base64.RawURLEncoding.EncodeToString(pub.Bytes()), nil
}
