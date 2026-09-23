package oauth2

import (
	"crypto/sha256"
	"crypto/subtle"
	"strings"

	"github.com/blocknextai/go-packages/base64"
)

type CodeChallengeMethod string

const (
	S256CodeChallengeMethod CodeChallengeMethod = "S256"
)

var (
	AllCodeChallengeMethods = map[CodeChallengeMethod]struct{}{
		S256CodeChallengeMethod: {},
	}
)

func (m CodeChallengeMethod) String() string {
	return string(m)
}

func (m CodeChallengeMethod) IsValid() bool {
	_, ok := AllCodeChallengeMethods[m]
	return ok
}

func (m CodeChallengeMethod) Verify(codeVerifier string, codeChallenge string) bool {
	if m != S256CodeChallengeMethod || strings.TrimSpace(codeVerifier) == "" {
		return false
	}

	hash := sha256.Sum256([]byte(codeVerifier))
	expected := base64.RawURLEncode(hash[:])

	return subtle.ConstantTimeCompare([]byte(expected), []byte(codeChallenge)) == 1
}
