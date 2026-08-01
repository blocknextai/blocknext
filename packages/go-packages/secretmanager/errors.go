package secretmanager

import (
	"errors"
)

var (
	// ErrCiphertextTooShort indicates the ciphertext is shorter than the
	// required GCM nonce.
	ErrCiphertextTooShort = errors.New("ciphertext too short")
	// ErrInvalidKeyFormat indicates the secret key is not valid hexadecimal.
	ErrInvalidKeyFormat = errors.New("invalid key format")
	// ErrInvalidKeyLength indicates the decoded secret key is not 32 bytes.
	ErrInvalidKeyLength = errors.New("invalid key length")
)
