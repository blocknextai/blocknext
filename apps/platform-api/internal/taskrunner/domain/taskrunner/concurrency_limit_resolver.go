package taskrunner

import (
	"context"

	"github.com/google/uuid"
)

type ConcurrencyLimitResolver interface {
	GetMaxConcurrentTasks(ctx context.Context, organizationID uuid.UUID) int64
}
