package credentials

import (
	"context"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type CredentialRepository interface {
	GetAllByOwner(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, searchQuery string, offset int, limit int) ([]*Credential, int64, error)
	GetAllByOwnerAndKeys(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, keys []string) ([]*Credential, error)
	GetByIDAndOwner(ctx context.Context, id uuid.UUID, ownerType commonDomain.OwnerType, ownerID uuid.UUID) (*Credential, error)
	Create(ctx context.Context, credential *Credential) error
	Update(ctx context.Context, credential *Credential) error
	Delete(ctx context.Context, credential *Credential) error
}
