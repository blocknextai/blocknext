package taskrunner

import (
	"context"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
)

type TaskLifecycleManager interface {
	StartTask(ctx context.Context, task *taskRunnerDomainTask.Task) error
	HandleTaskSuccess(ctx context.Context, task *taskRunnerDomainTask.Task) error
	HandleTaskFailure(ctx context.Context, task *taskRunnerDomainTask.Task, err error) error
	HandleTaskCancellation(ctx context.Context, task *taskRunnerDomainTask.Task) error
	CreateNodeExecutions(ctx context.Context, task *taskRunnerDomainTask.Task) error
}
