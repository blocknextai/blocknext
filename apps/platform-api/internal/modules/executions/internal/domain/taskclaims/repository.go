package taskclaims

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskClaimRepository interface {
	Create(ctx context.Context, taskClaim *TaskClaim) error
	GetByTaskExecutionID(ctx context.Context, taskExecutionID uuid.UUID) (*TaskClaim, error)
	GetByTaskExecutionIDForUpdate(ctx context.Context, taskExecutionID uuid.UUID) (*TaskClaim, error)
	GetAllStale(ctx context.Context, staleAfter time.Duration) ([]*TaskClaim, error)
	Update(ctx context.Context, taskClaim *TaskClaim) error
}
