package auth

import (
	"context"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
)

type AuthenticatedAPIKey[S ~string] struct {
	ID        uuid.UUID
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	Scopes    []S
}

type APIKeyValidator[S ~string] interface {
	Validate(ctx context.Context, rawKey string) (*AuthenticatedAPIKey[S], error)
}
