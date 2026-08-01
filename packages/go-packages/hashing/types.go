// Package hashing defines the Hasher interface for password hashing and
// verification.
package hashing

// Hasher generates and verifies password hashes.
type Hasher interface {
	// Generate returns a hash of the given password.
	Generate(password string) (string, error)
	// Compare reports whether the password matches the given hash.
	Compare(password, hash string) bool
}
