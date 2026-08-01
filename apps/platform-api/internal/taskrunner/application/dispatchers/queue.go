package dispatchers

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	executionsApplicationTaskClaims "github.com/blocknextai/platform-api/internal/executions/application/taskclaims"
	taskRunnerApplicationTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/application/taskrunner"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/blocknextai/platform-api/internal/taskrunner/domain/taskqueue"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
)

const (
	reclaimInterval = 30 * time.Second
)

var (
	ErrTaskExceededMaxRetries = errors.New("task exceeded max retries")
)

type Queue struct {
	queue            taskqueue.TaskQueue
	workerPool       taskRunnerDomainTaskRunner.WorkerPool
	executor         *taskRunnerApplicationTaskRunner.TaskExecutor
	taskClaimService executionsApplicationTaskClaims.TaskClaimService
	idleTimeout      time.Duration
	maxRetries       int
	wg               sync.WaitGroup
	startOnce        sync.Once
	stopOnce         sync.Once
	cancelCtx        context.CancelFunc
}

func NewQueue(
	queue taskqueue.TaskQueue,
	workerPool taskRunnerDomainTaskRunner.WorkerPool,
	executor *taskRunnerApplicationTaskRunner.TaskExecutor,
	taskClaimService executionsApplicationTaskClaims.TaskClaimService,
	idleTimeout time.Duration,
	maxRetries int,
) taskRunnerDomainTaskRunner.TaskDispatcher {
	return &Queue{
		queue:            queue,
		workerPool:       workerPool,
		executor:         executor,
		taskClaimService: taskClaimService,
		idleTimeout:      idleTimeout,
		maxRetries:       maxRetries,
	}
}

func (d *Queue) Ping(ctx context.Context) error {
	return d.queue.Ping(ctx)
}

func (d *Queue) Dispatch(ctx context.Context, envelope taskRunnerDomainTask.TaskEnvelope) error {
	return d.queue.Push(ctx, envelope)
}

func (d *Queue) Start(ctx context.Context) error {
	var startErr error
	d.startOnce.Do(func() {
		runCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
		d.cancelCtx = cancel

		messages, err := d.queue.Subscribe(runCtx)
		if err != nil {
			cancel()
			startErr = err
			return
		}

		d.wg.Add(2)
		go d.consume(runCtx, messages)
		go d.reclaimLoop(runCtx)
	})
	return startErr
}

func (d *Queue) Stop() error {
	d.stopOnce.Do(func() {
		if d.cancelCtx != nil {
			d.cancelCtx()
		}
		d.wg.Wait()
	})
	return nil
}

func (d *Queue) consume(ctx context.Context, messages <-chan taskqueue.TaskMessage) {
	defer d.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-messages:
			if !ok {
				return
			}
			d.handleMessage(ctx, msg)
		}
	}
}

func (d *Queue) handleMessage(ctx context.Context, msg taskqueue.TaskMessage) {
	envelope := msg.Envelope
	msgID := msg.ID

	err := d.workerPool.SubmitWithContext(ctx, func() {
		execErr := d.executor.Execute(ctx, envelope)
		if execErr == nil {
			if resetErr := d.taskClaimService.ResetRetryCount(ctx, envelope.TaskID); resetErr != nil {
				slog.WarnContext(ctx, "failed to reset retry count after successful execution",
					"component", "queue_dispatcher",
					"task_id", envelope.TaskID,
					"error", resetErr)
			}
			if ackErr := d.queue.Ack(ctx, msgID); ackErr != nil {
				slog.ErrorContext(ctx, "failed to ack task message after execution",
					"component", "queue_dispatcher",
					"task_id", envelope.TaskID,
					"message_id", msgID,
					"error", ackErr)
			}
			return
		}

		if errors.Is(execErr, taskRunnerApplicationTaskRunner.ErrTaskClaimNotAvailable) {
			slog.InfoContext(ctx, "leaving message for redelivery, claim held by another worker",
				"component", "queue_dispatcher",
				"task_id", envelope.TaskID,
				"message_id", msgID)
			return
		}

		newCount, incErr := d.taskClaimService.IncrementRetryCount(ctx, envelope.TaskID)
		if incErr != nil {
			slog.ErrorContext(ctx, "failed to increment retry count, leaving message for autoclaim",
				"component", "queue_dispatcher",
				"task_id", envelope.TaskID,
				"message_id", msgID,
				"error", incErr)
			return
		}

		if newCount >= d.maxRetries {
			slog.ErrorContext(ctx, "task exceeded max retries, marking as failed",
				"component", "queue_dispatcher",
				"task_id", envelope.TaskID,
				"message_id", msgID,
				"retries", newCount,
				"max_retries", d.maxRetries,
				"last_error", execErr)
			poisonErr := errors.Join(ErrTaskExceededMaxRetries, execErr)
			if markErr := d.executor.MarkAsFailed(ctx, envelope, poisonErr); markErr != nil {
				slog.ErrorContext(ctx, "failed to mark poisoned task as failed",
					"component", "queue_dispatcher",
					"task_id", envelope.TaskID,
					"error", markErr)
			}
			if ackErr := d.queue.Ack(ctx, msgID); ackErr != nil {
				slog.ErrorContext(ctx, "failed to ack poisoned task message",
					"component", "queue_dispatcher",
					"task_id", envelope.TaskID,
					"message_id", msgID,
					"error", ackErr)
			}
			return
		}

		slog.WarnContext(ctx, "task execution failed, leaving message for redelivery",
			"component", "queue_dispatcher",
			"task_id", envelope.TaskID,
			"message_id", msgID,
			"retries", newCount,
			"max_retries", d.maxRetries,
			"error", execErr)
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to submit task to worker pool, leaving message for autoclaim",
			"component", "queue_dispatcher",
			"task_id", envelope.TaskID,
			"message_id", msgID,
			"error", err)
	}
}

func (d *Queue) reclaimLoop(ctx context.Context) {
	defer d.wg.Done()

	ticker := time.NewTicker(reclaimInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			reclaimed, err := d.queue.ReclaimIdle(ctx, d.idleTimeout)
			if err != nil {
				slog.ErrorContext(ctx, "failed to reclaim idle messages",
					"component", "queue_dispatcher",
					"error", err)
				continue
			}
			if reclaimed > 0 {
				slog.InfoContext(ctx, "reclaimed idle messages from pending entry list",
					"component", "queue_dispatcher",
					"count", reclaimed)
			}
		}
	}
}
