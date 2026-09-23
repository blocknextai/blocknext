package organizationusers

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrganizationUserRepository interface {
	GetAllByIDs(ctx context.Context, ids []uuid.UUID) ([]*OrganizationUser, error)
	GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, searchQuery string, offset int, limit int) ([]*OrganizationUser, int64, error)
	GetAllByOrganizationIDAndUserIDs(ctx context.Context, organizationID uuid.UUID, userIDs []uuid.UUID) ([]*OrganizationUser, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*OrganizationUser, error)
	GetByID(ctx context.Context, id uuid.UUID) (*OrganizationUser, error)
	GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*OrganizationUser, error)
	GetByOrganizationIDAndUserID(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (*OrganizationUser, error)
	CountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error)
	HasOwner(ctx context.Context, organizationID uuid.UUID) (bool, error)
	Create(ctx context.Context, organizationUser *OrganizationUser) error
	Update(ctx context.Context, organizationUser *OrganizationUser) error
	Delete(ctx context.Context, organizationUser *OrganizationUser) error
	DeleteAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, deletedAt time.Time) error
}
