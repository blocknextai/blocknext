// Package hex wraps the standard encoding/hex package for hexadecimal encoding
// and decoding.
package hex

import (
	stdhex "encoding/hex"
)

// Encode returns the hexadecimal encoding of data.
func Encode(data []byte) string {
	return stdhex.EncodeToString(data)
}

// Decode returns the bytes represented by the hexadecimal string data.
func Decode(data string) ([]byte, error) {
	return stdhex.DecodeString(data)
}
