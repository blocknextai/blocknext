package taskexecutions

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskExecutionRepository interface {
	GetAllByOrganizationID(
		ctx context.Context,
		organizationID uuid.UUID,
		searchQuery string,
		offset int,
		limit int,
	) ([]*TaskExecution, int64, error)
	GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*TaskExecution, error)
	GetByID(ctx context.Context, id uuid.UUID) (*TaskExecution, error)
	GetAllByStatuses(ctx context.Context, statuses []string) ([]*TaskExecution, error)

	Create(ctx context.Context, taskExecution *TaskExecution) error
	Update(ctx context.Context, taskExecution *TaskExecution) error
	Delete(ctx context.Context, taskExecution *TaskExecution) error
	BulkDelete(ctx context.Context, ids []uuid.UUID, organizationID uuid.UUID, now time.Time) error
	GetAllByIDsAndOrganizationID(ctx context.Context, ids []uuid.UUID, organizationID uuid.UUID) ([]*TaskExecution, error)
}
