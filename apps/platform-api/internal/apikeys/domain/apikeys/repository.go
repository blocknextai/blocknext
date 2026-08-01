package apikeys

import (
	"context"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type APIKeyRepository interface {
	GetByOwnerAndID(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, id uuid.UUID) (*APIKey, error)
	GetByKeyHash(ctx context.Context, keyHash string) (*APIKey, error)
	GetAllByOwner(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, searchQuery string, offset int, limit int) ([]*APIKey, int64, error)
	Create(ctx context.Context, apiKey *APIKey) error
	Update(ctx context.Context, apiKey *APIKey) error
	Delete(ctx context.Context, apiKey *APIKey) error
}
