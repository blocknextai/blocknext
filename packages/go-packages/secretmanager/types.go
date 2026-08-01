// Package secretmanager provides authenticated symmetric encryption of
// arbitrary values using AES-256-GCM, with results serialized to JSON and
// encoded as base64.
package secretmanager

// SecretManager encrypts and decrypts arbitrary values using a symmetric key.
type SecretManager interface {
	// Encrypt serializes plaintext to JSON, encrypts it, and returns the
	// base64-encoded ciphertext.
	Encrypt(plaintext any) (string, error)
	// Decrypt decodes and decrypts encryptedText and unmarshals the result
	// into output.
	Decrypt(encryptedText string, output any) error
}
