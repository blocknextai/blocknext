package taskrunner

import (
	"context"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
)

type TaskExecutionCoordinator interface {
	ExecuteTask(ctx context.Context, task *taskRunnerDomainTask.Task) error
	CancelExecution(ctx context.Context, taskID uuid.UUID) error
}
