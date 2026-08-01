package organizations

import (
	"context"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	Create(ctx context.Context, organization *Organization) error
	Update(ctx context.Context, organization *Organization) error
	Delete(ctx context.Context, organization *Organization) error
	GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*Organization, error)
}
