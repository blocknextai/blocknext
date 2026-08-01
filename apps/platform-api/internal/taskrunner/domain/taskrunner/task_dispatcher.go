package taskrunner

import (
	"context"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
)

type TaskDispatcher interface {
	Ping(ctx context.Context) error
	Dispatch(ctx context.Context, envelope taskRunnerDomainTask.TaskEnvelope) error
	Start(ctx context.Context) error
	Stop() error
}
