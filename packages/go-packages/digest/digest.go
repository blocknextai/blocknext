// Package digest provides convenience helpers for computing cryptographic
// digests.
package digest

import (
	"crypto/sha256"

	"github.com/blocknextai/go-packages/hex"
)

// SHA256Hex returns the hexadecimal-encoded SHA-256 digest of data.
func SHA256Hex(data string) string {
	sum := sha256.Sum256([]byte(data))
	return hex.Encode(sum[:])
}
