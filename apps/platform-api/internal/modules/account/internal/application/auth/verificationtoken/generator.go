package verificationtoken

import (
	"crypto/rand"

	"github.com/blocknextai/go-packages/digest"
	"github.com/blocknextai/go-packages/hex"
)

const (
	tokenByteLength = 32
)

func Generate() (plainToken string, tokenHash string, err error) {
	buf := make([]byte, tokenByteLength)
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}

	plainToken = hex.Encode(buf)
	tokenHash = Hash(plainToken)
	return plainToken, tokenHash, nil
}

func Hash(plainToken string) string {
	return digest.SHA256Hex(plainToken)
}
