package apikeys

import (
	"crypto/rand"
	"strings"

	"github.com/blocknextai/go-packages/digest"
	"github.com/blocknextai/go-packages/hex"
)

const (
	keyByteLength   = 32
	keyFormatPrefix = "bnx_"
)

type GeneratedKey struct {
	Plain string
	Hash  string
}

func GenerateKey() (*GeneratedKey, error) {
	bytes := make([]byte, keyByteLength)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	var builder strings.Builder
	builder.WriteString(keyFormatPrefix)
	builder.WriteString(hex.Encode(bytes))
	plain := builder.String()

	return &GeneratedKey{
		Plain: plain,
		Hash:  digest.SHA256Hex(plain),
	}, nil
}
