package grants

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
)

type Grant struct {
	database.BaseEntity

	ClientID       string
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Scopes         mcpOAuthDomainOAuth2.Scopes
	Resource       string
	RevokedAt      *time.Time
}

func New(
	clientID string,
	userID uuid.UUID,
	organizationID uuid.UUID,
	scopes mcpOAuthDomainOAuth2.Scopes,
	resource string,
) (*Grant, error) {
	now := time.Now().UTC()

	grant := &Grant{
		ID:             bnuuid.NewV7(),
		CreatedAt:      now,
		UpdatedAt:      now,
		ClientID:       clientID,
		UserID:         userID,
		OrganizationID: organizationID,
		Scopes:         scopes,
		Resource:       resource,
	}

	return grant.validateThenReturn()
}

func (g *Grant) IsActive() bool {
	return g.RevokedAt == nil && g.DeletedAt == nil
}

func (g *Grant) Update(scopes mcpOAuthDomainOAuth2.Scopes) (*Grant, error) {
	g.Scopes = scopes
	g.UpdatedAt = time.Now().UTC()

	return g.validateThenReturn()
}

func (g *Grant) Revoke() (*Grant, error) {
	now := time.Now().UTC()

	g.RevokedAt = new(now)
	g.UpdatedAt = now

	return g.validateThenReturn()
}

func (g *Grant) validateThenReturn() (*Grant, error) {
	if strings.TrimSpace(g.ClientID) == "" {
		return nil, ErrInvalidClientID
	}

	if g.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if g.OrganizationID == uuid.Nil {
		return nil, ErrInvalidOrganizationID
	}

	if strings.TrimSpace(g.Resource) == "" {
		return nil, ErrInvalidResource
	}

	if len(g.Scopes) == 0 || !g.Scopes.IsValid() {
		return nil, ErrInvalidScopes
	}

	return g, nil
}
