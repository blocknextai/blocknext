package application

import (
	"crypto/rand"

	"github.com/blocknextai/go-packages/digest"
	"github.com/blocknextai/go-packages/hex"
)

func GenerateWebhookToken() (plain string, hash string, err error) {
	bytes := make([]byte, 32)
	if _, err = rand.Read(bytes); err != nil {
		return "", "", err
	}
	plain = hex.Encode(bytes)
	hash = digest.SHA256Hex(plain)
	return plain, hash, nil
}

func HashWebhookToken(plain string) string {
	return digest.SHA256Hex(plain)
}
