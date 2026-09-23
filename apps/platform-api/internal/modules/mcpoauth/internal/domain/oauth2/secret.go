package oauth2

import (
	"crypto/rand"

	"github.com/blocknextai/go-packages/digest"
	"github.com/blocknextai/go-packages/hex"
)

const (
	secretByteLength = 32
)

type Secret struct {
	Plain string
	Hash  string
}

func GenerateSecret(prefix string) (*Secret, error) {
	bytes := make([]byte, secretByteLength)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	plain := prefix + hex.Encode(bytes)

	return &Secret{
		Plain: plain,
		Hash:  HashSecret(plain),
	}, nil
}

func HashSecret(plain string) string {
	return digest.SHA256Hex(plain)
}
