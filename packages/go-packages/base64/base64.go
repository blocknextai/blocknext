// Package base64 wraps the standard encoding/base64 package with standard and
// URL-safe encoding and decoding helpers.
package base64

import (
	"bytes"
	stdbase64 "encoding/base64"
	"io"
	"strings"
)

// Encode returns the standard base64 encoding of data.
func Encode(data []byte) string {
	return stdbase64.StdEncoding.EncodeToString(data)
}

// Decode returns the bytes represented by the standard base64 string data.
func Decode(data string) ([]byte, error) {
	reader := stdbase64.NewDecoder(stdbase64.StdEncoding, strings.NewReader(data))
	return io.ReadAll(reader)
}

// DecodeReader returns an io.Reader that decodes the standard base64 string data.
func DecodeReader(data string) io.Reader {
	return stdbase64.NewDecoder(stdbase64.StdEncoding, strings.NewReader(data))
}

// EncodeReader reads all data from reader and returns its standard base64 encoding.
func EncodeReader(reader io.Reader) (string, error) {
	var buf bytes.Buffer
	encoder := stdbase64.NewEncoder(stdbase64.StdEncoding, &buf)

	if _, err := io.Copy(encoder, reader); err != nil {
		return "", err
	}

	if err := encoder.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// URLEncode returns the URL-safe base64 encoding of data.
func URLEncode(data []byte) string {
	return stdbase64.URLEncoding.EncodeToString(data)
}

// URLDecode returns the bytes represented by the URL-safe base64 string data.
func URLDecode(data string) ([]byte, error) {
	return stdbase64.URLEncoding.DecodeString(data)
}

// RawURLEncode returns the unpadded URL-safe base64 encoding of data.
func RawURLEncode(data []byte) string {
	return stdbase64.RawURLEncoding.EncodeToString(data)
}

// RawURLDecode returns the bytes represented by the unpadded URL-safe base64 string data.
func RawURLDecode(data string) ([]byte, error) {
	return stdbase64.RawURLEncoding.DecodeString(data)
}
