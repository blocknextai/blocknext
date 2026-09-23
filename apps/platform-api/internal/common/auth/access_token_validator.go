package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthenticatedAccessToken struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Scopes         []string
}

type AccessTokenValidator interface {
	Validate(ctx context.Context, rawToken string, resourcePath string) (*AuthenticatedAccessToken, error)
}
