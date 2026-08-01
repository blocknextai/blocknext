package dispatchers

import (
	"context"

	taskRunnerApplicationTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/application/taskrunner"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
)

type Embedded struct {
	workerPool taskRunnerDomainTaskRunner.WorkerPool
	executor   *taskRunnerApplicationTaskRunner.TaskExecutor
}

func NewEmbedded(
	workerPool taskRunnerDomainTaskRunner.WorkerPool,
	executor *taskRunnerApplicationTaskRunner.TaskExecutor,
) taskRunnerDomainTaskRunner.TaskDispatcher {
	return &Embedded{
		workerPool: workerPool,
		executor:   executor,
	}
}

func (d *Embedded) Ping(ctx context.Context) error {
	return nil
}

func (d *Embedded) Dispatch(ctx context.Context, envelope taskRunnerDomainTask.TaskEnvelope) error {
	detachedCtx := context.WithoutCancel(ctx)
	return d.workerPool.SubmitWithContext(detachedCtx, func() {
		_ = d.executor.Execute(detachedCtx, envelope)
	})
}

func (d *Embedded) Start(ctx context.Context) error {
	return nil
}

func (d *Embedded) Stop() error {
	return nil
}
