package authorizationcodes

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	CodePrefix = "mcpa_"
)

type AuthorizationCode struct {
	database.BaseEntity

	GrantID             uuid.UUID
	ClientID            string
	CodeHash            string
	RedirectURI         string
	CodeChallenge       string
	CodeChallengeMethod mcpOAuthDomainOAuth2.CodeChallengeMethod
	Resource            string
	Scopes              mcpOAuthDomainOAuth2.Scopes
	ExpiresAt           time.Time
	UsedAt              *time.Time
}

func New(
	grantID uuid.UUID,
	clientID string,
	codeHash string,
	redirectURI string,
	codeChallenge string,
	codeChallengeMethod mcpOAuthDomainOAuth2.CodeChallengeMethod,
	resource string,
	scopes mcpOAuthDomainOAuth2.Scopes,
	ttl time.Duration,
) (*AuthorizationCode, error) {
	now := time.Now().UTC()

	code := &AuthorizationCode{
		ID:                  bnuuid.NewV7(),
		CreatedAt:           now,
		UpdatedAt:           now,
		GrantID:             grantID,
		ClientID:            clientID,
		CodeHash:            codeHash,
		RedirectURI:         redirectURI,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		Resource:            resource,
		Scopes:              scopes,
		ExpiresAt:           now.Add(ttl),
	}

	return code.validateThenReturn()
}

func (c *AuthorizationCode) IsExpired() bool {
	return time.Now().UTC().After(c.ExpiresAt)
}

func (c *AuthorizationCode) IsUsed() bool {
	return c.UsedAt != nil
}

func (c *AuthorizationCode) VerifyCodeVerifier(codeVerifier string) bool {
	return c.CodeChallengeMethod.Verify(codeVerifier, c.CodeChallenge)
}

func (c *AuthorizationCode) Use() (*AuthorizationCode, error) {
	now := time.Now().UTC()

	c.UsedAt = new(now)
	c.UpdatedAt = now

	return c.validateThenReturn()
}

func (c *AuthorizationCode) validateThenReturn() (*AuthorizationCode, error) {
	if c.GrantID == uuid.Nil {
		return nil, ErrInvalidGrantID
	}

	if strings.TrimSpace(c.ClientID) == "" {
		return nil, ErrInvalidClientID
	}

	if strings.TrimSpace(c.CodeHash) == "" {
		return nil, ErrInvalidCodeHash
	}

	if strings.TrimSpace(c.RedirectURI) == "" {
		return nil, ErrInvalidRedirectURI
	}

	if strings.TrimSpace(c.CodeChallenge) == "" || !c.CodeChallengeMethod.IsValid() {
		return nil, ErrInvalidCodeChallenge
	}

	if strings.TrimSpace(c.Resource) == "" {
		return nil, ErrInvalidResource
	}

	if len(c.Scopes) == 0 || !c.Scopes.IsValid() {
		return nil, ErrInvalidScopes
	}

	return c, nil
}
