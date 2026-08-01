// Package pkce generates PKCE (Proof Key for Code Exchange) code verifiers and
// challenges for OAuth 2.0 authorization flows.
package pkce

import (
	"crypto/rand"
	"crypto/sha256"

	"github.com/blocknextai/go-packages/base64"
)

// CodeVerifierLength is the number of random bytes used to build a code verifier.
const CodeVerifierLength = 32

// Challenge holds a PKCE code verifier together with its derived code challenge
// and the challenge method.
type Challenge struct {
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

// GenerateChallenge returns a new PKCE Challenge using the S256 method.
func GenerateChallenge() (*Challenge, error) {
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, err
	}

	codeChallenge := generateCodeChallenge(codeVerifier)

	return &Challenge{
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: "S256",
	}, nil
}

func generateCodeVerifier() (string, error) {
	bytes := make([]byte, CodeVerifierLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncode(bytes), nil
}

func generateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncode(hash[:])
}
