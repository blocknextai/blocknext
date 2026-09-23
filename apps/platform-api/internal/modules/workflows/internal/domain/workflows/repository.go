package workflows

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowRepository interface {
	GetAllByOrganizationID(ctx context.Context, organizationID uuid.UUID, searchQuery string, offset int, limit int) ([]*Workflow, int64, error)
	GetByOrganizationIDAndID(ctx context.Context, organizationID uuid.UUID, workflowID uuid.UUID) (*Workflow, error)
	GetCountByOrganizationID(ctx context.Context, organizationID uuid.UUID) (int64, error)
	Create(ctx context.Context, workflow *Workflow) error
	Update(ctx context.Context, workflow *Workflow) error
	Delete(ctx context.Context, workflow *Workflow) error
}
