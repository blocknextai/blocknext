package taskqueue

import (
	"context"
	"time"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
)

type TaskMessage struct {
	ID       string
	Envelope taskRunnerDomainTask.TaskEnvelope
}

type TaskQueue interface {
	Ping(ctx context.Context) error
	Push(ctx context.Context, envelope taskRunnerDomainTask.TaskEnvelope) error
	Subscribe(ctx context.Context) (<-chan TaskMessage, error)
	Ack(ctx context.Context, msgID string) error
	ReclaimIdle(ctx context.Context, idleTimeout time.Duration) (int, error)
	Close() error
}
