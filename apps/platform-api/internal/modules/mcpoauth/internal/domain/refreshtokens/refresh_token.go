package refreshtokens

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

const (
	TokenPrefix = "mcpr_"
)

type RefreshToken struct {
	database.BaseEntity

	GrantID   uuid.UUID
	ClientID  string
	TokenHash string
	Scopes    mcpOAuthDomainOAuth2.Scopes
	Resource  string
	ExpiresAt time.Time
	UsedAt    *time.Time
	RevokedAt *time.Time
}

func New(
	grantID uuid.UUID,
	clientID string,
	tokenHash string,
	scopes mcpOAuthDomainOAuth2.Scopes,
	resource string,
	ttl time.Duration,
) (*RefreshToken, error) {
	now := time.Now().UTC()

	token := &RefreshToken{
		ID:        bnuuid.NewV7(),
		CreatedAt: now,
		UpdatedAt: now,
		GrantID:   grantID,
		ClientID:  clientID,
		TokenHash: tokenHash,
		Scopes:    scopes,
		Resource:  resource,
		ExpiresAt: now.Add(ttl),
	}

	return token.validateThenReturn()
}

func (t *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(t.ExpiresAt)
}

func (t *RefreshToken) IsUsed() bool {
	return t.UsedAt != nil
}

func (t *RefreshToken) IsRevoked() bool {
	return t.RevokedAt != nil
}

func (t *RefreshToken) Use() (*RefreshToken, error) {
	now := time.Now().UTC()

	t.UsedAt = new(now)
	t.UpdatedAt = now

	return t.validateThenReturn()
}

func (t *RefreshToken) validateThenReturn() (*RefreshToken, error) {
	if t.GrantID == uuid.Nil {
		return nil, ErrInvalidGrantID
	}

	if strings.TrimSpace(t.ClientID) == "" {
		return nil, ErrInvalidClientID
	}

	if strings.TrimSpace(t.TokenHash) == "" {
		return nil, ErrInvalidTokenHash
	}

	if strings.TrimSpace(t.Resource) == "" {
		return nil, ErrInvalidResource
	}

	if len(t.Scopes) == 0 || !t.Scopes.IsValid() {
		return nil, ErrInvalidScopes
	}

	return t, nil
}
