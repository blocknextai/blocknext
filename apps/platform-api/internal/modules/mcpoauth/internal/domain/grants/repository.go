package grants

import (
	"context"

	"github.com/google/uuid"
)

type GrantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Grant, error)
	GetActiveByClientIDAndUserIDAndOrganizationIDAndResource(ctx context.Context, clientID string, userID uuid.UUID, organizationID uuid.UUID, resource string) (*Grant, error)
	GetAllActiveByUserID(ctx context.Context, userID uuid.UUID, searchQuery string, offset int, limit int) ([]*Grant, int64, error)
	Create(ctx context.Context, grant *Grant) error
	Update(ctx context.Context, grant *Grant) error
}
