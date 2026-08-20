package usernonces

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/uuid"
	"github.com/blocknextai/platform-api/internal/account/domain"
)

type UserNonce struct {
	database.BaseEntity

	AuthProvider        domain.AuthProvider
	ProviderID          *string
	Nonce               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
}

func NewUserNonce(
	authProvider domain.AuthProvider,
	providerId *string,
	nonce string,
	codeVerifier string,
	codeChallenge string,
	codeChallengeMethod string,
) (*UserNonce, error) {
	utcNow := time.Now().UTC()

	userNonce := &UserNonce{
		ID:                  uuid.NewV7(),
		CreatedAt:           utcNow,
		UpdatedAt:           utcNow,
		DeletedAt:           nil,
		AuthProvider:        authProvider,
		ProviderID:          providerId,
		Nonce:               nonce,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:           time.Now().UTC().Add(time.Minute * 5),
	}

	return userNonce.validateThenReturn()
}

func (u *UserNonce) Update(nonce string, codeVerifier string, codeChallenge string, codeChallengeMethod string, expiresAt time.Time) (*UserNonce, error) {
	u.UpdatedAt = time.Now().UTC()

	u.Nonce = nonce
	u.CodeVerifier = codeVerifier
	u.CodeChallenge = codeChallenge
	u.CodeChallengeMethod = codeChallengeMethod
	u.ExpiresAt = expiresAt

	return u.validateThenReturn()
}

func (u *UserNonce) Delete() (*UserNonce, error) {
	utcNow := time.Now().UTC()

	u.UpdatedAt = utcNow
	u.DeletedAt = new(utcNow)
	return u.validateThenReturn()
}

func (u *UserNonce) validateThenReturn() (*UserNonce, error) {
	if strings.TrimSpace(u.Nonce) == "" {
		return nil, ErrNonceIsRequired
	}

	if strings.TrimSpace(u.CodeVerifier) == "" {
		return nil, ErrCodeVerifierIsRequired
	}

	if strings.TrimSpace(u.CodeChallenge) == "" {
		return nil, ErrCodeChallengeIsRequired
	}

	if strings.TrimSpace(u.CodeChallengeMethod) == "" {
		return nil, ErrCodeChallengeMethodIsRequired
	}

	if u.ExpiresAt.IsZero() {
		return nil, ErrExpiresAtIsRequired
	}

	if u.ExpiresAt.Before(time.Now().UTC()) {
		return nil, ErrExpiresAtInThePast
	}

	return u, nil
}
