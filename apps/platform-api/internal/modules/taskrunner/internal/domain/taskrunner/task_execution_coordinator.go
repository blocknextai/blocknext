package taskrunner

import (
	"context"

	"github.com/google/uuid"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/task"
)

type TaskExecutionCoordinator interface {
	ExecuteTask(ctx context.Context, task *taskRunnerDomainTask.Task) error
	CancelExecution(ctx context.Context, taskID uuid.UUID) error
}
