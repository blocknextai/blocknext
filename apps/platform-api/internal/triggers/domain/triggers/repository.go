package triggers

import (
	"context"

	"github.com/google/uuid"
)

type TriggerRepository interface {
	GetAllActive(ctx context.Context) ([]*Trigger, error)
	GetAllByOrganizationID(
		ctx context.Context,
		organizationID uuid.UUID,
		searchQuery string,
		offset int,
		limit int,
	) ([]*Trigger, int64, error)
	GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*Trigger, error)
	GetByWebhookTokenHash(ctx context.Context, tokenHash string) (*Trigger, error)
	Create(ctx context.Context, trigger *Trigger) error
	Update(ctx context.Context, trigger *Trigger) error
	Delete(ctx context.Context, trigger *Trigger) error
}
