package taskrunner

import (
	"context"

	"github.com/google/uuid"
)

type SemaphoreManager interface {
	Ping(ctx context.Context) error
	AcquireSemaphore(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID, maxConcurrentExecutions int64) (chan struct{}, error)
	ReleaseSemaphore(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID, semaphore chan struct{}) error
	HeartbeatSemaphore(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID) error
}
