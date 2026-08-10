package taskrunner

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/blocknextai/go-packages/dag"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	executionsApplicationTaskClaims "github.com/blocknextai/platform-api/internal/executions/application/taskclaims"
	executionsApplicationTaskexecutions "github.com/blocknextai/platform-api/internal/executions/application/taskexecutions"
	executionsDomainTaskExecutions "github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/google/uuid"
)

type TaskExecutor struct {
	taskExecutionService     executionsApplicationTaskexecutions.TaskExecutionService
	taskClaimService         executionsApplicationTaskClaims.TaskClaimService
	nodeExecutionService     executionsApplicationNodeExecutions.NodeExecutionService
	executionCoordinator     taskRunnerDomainTaskRunner.TaskExecutionCoordinator
	lifecycleManager         taskRunnerDomainTaskRunner.TaskLifecycleManager
	semaphoreManager         taskRunnerDomainTaskRunner.SemaphoreManager
	contextManager           taskRunnerDomainTaskRunner.ContextManager
	concurrencyLimitResolver taskRunnerDomainTaskRunner.ConcurrencyLimitResolver
	workerID                 string
	heartbeatInterval        time.Duration
	maxExecutionTime         time.Duration
}

func NewTaskExecutor(
	taskExecutionService executionsApplicationTaskexecutions.TaskExecutionService,
	taskClaimService executionsApplicationTaskClaims.TaskClaimService,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	executionCoordinator taskRunnerDomainTaskRunner.TaskExecutionCoordinator,
	lifecycleManager taskRunnerDomainTaskRunner.TaskLifecycleManager,
	semaphoreManager taskRunnerDomainTaskRunner.SemaphoreManager,
	contextManager taskRunnerDomainTaskRunner.ContextManager,
	concurrencyLimitResolver taskRunnerDomainTaskRunner.ConcurrencyLimitResolver,
	workerID string,
	heartbeatInterval time.Duration,
	maxExecutionTime time.Duration,
) *TaskExecutor {
	return &TaskExecutor{
		taskExecutionService:     taskExecutionService,
		taskClaimService:         taskClaimService,
		nodeExecutionService:     nodeExecutionService,
		executionCoordinator:     executionCoordinator,
		lifecycleManager:         lifecycleManager,
		semaphoreManager:         semaphoreManager,
		contextManager:           contextManager,
		concurrencyLimitResolver: concurrencyLimitResolver,
		workerID:                 workerID,
		heartbeatInterval:        heartbeatInterval,
		maxExecutionTime:         maxExecutionTime,
	}
}

func runPeriodic(ctx context.Context, interval time.Duration, tick func() (stop bool)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if tick() {
				return
			}
		}
	}
}

func (e *TaskExecutor) runHeartbeat(ctx context.Context, taskID uuid.UUID) {
	runPeriodic(ctx, e.heartbeatInterval, func() (stop bool) {
		taskExecution, err := e.taskExecutionService.GetByID(ctx, taskID)
		if err == nil && !isActiveStatus(taskExecution.Status) {
			slog.InfoContext(ctx, "task status changed externally, cancelling local execution",
				"component", "task_executor",
				"task_id", taskID,
				"status", taskExecution.Status)
			e.contextManager.CancelContext(taskID)
			return true
		}

		if err := e.taskClaimService.Heartbeat(ctx, taskID, e.workerID); err != nil {
			slog.WarnContext(ctx, "failed to heartbeat task claim",
				"component", "task_executor",
				"task_id", taskID,
				"worker_id", e.workerID,
				"error", err)
		}

		return false
	})
}

func (e *TaskExecutor) MarkAsFailed(ctx context.Context, envelope taskRunnerDomainTask.TaskEnvelope, reason error) error {
	taskExecution, err := e.taskExecutionService.GetByID(ctx, envelope.TaskID)
	if err != nil {
		return err
	}

	task := &taskRunnerDomainTask.Task{
		ID:                 taskExecution.ID,
		OrganizationID:     taskExecution.OrganizationID,
		TriggeredByUserID:  taskExecution.TriggeredByUserID,
		ExecutionContext:   taskExecution.ExecutionContext,
		ContextItemID:      taskExecution.ContextItemID,
		StartTime:          taskExecution.StartedAt,
		NodeExecutionIDMap: make(map[string]uuid.UUID),
	}
	if task.StartTime == nil {
		now := time.Now().UTC()
		task.StartTime = new(now)
	}

	return e.lifecycleManager.HandleTaskFailure(ctx, task, reason)
}

func (e *TaskExecutor) Execute(ctx context.Context, envelope taskRunnerDomainTask.TaskEnvelope) error {
	claimed, err := e.taskClaimService.ClaimByID(ctx, envelope.TaskID, e.workerID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to claim task",
			"component", "task_executor",
			"task_id", envelope.TaskID,
			"worker_id", e.workerID,
			"error", err)
		return err
	}
	if !claimed {
		taskExecution, loadErr := e.taskExecutionService.GetByID(ctx, envelope.TaskID)
		if loadErr == nil && !isActiveStatus(taskExecution.Status) {
			slog.InfoContext(ctx, "task not in active state, acking stale message",
				"component", "task_executor",
				"task_id", envelope.TaskID,
				"status", taskExecution.Status)
			return nil
		}
		slog.InfoContext(ctx, "task claim held by another worker, skipping",
			"component", "task_executor",
			"task_id", envelope.TaskID)
		return ErrTaskClaimNotAvailable
	}

	heartbeatCtx, stopHeartbeat := context.WithCancel(ctx)
	defer stopHeartbeat()
	go e.runHeartbeat(heartbeatCtx, envelope.TaskID)

	claimReleased := false
	defer func() {
		if claimReleased {
			return
		}
		if releaseErr := e.taskClaimService.ReleaseClaim(ctx, envelope.TaskID, e.workerID); releaseErr != nil {
			slog.ErrorContext(ctx, "failed to release task claim",
				"component", "task_executor",
				"task_id", envelope.TaskID,
				"worker_id", e.workerID,
				"error", releaseErr)
		}
	}()

	taskExecution, err := e.taskExecutionService.GetByID(ctx, envelope.TaskID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to load task execution",
			"component", "task_executor",
			"task_id", envelope.TaskID,
			"error", err)
		return err
	}

	if !isActiveStatus(taskExecution.Status) {
		slog.InfoContext(ctx, "task not in active state after claim, skipping execution",
			"component", "task_executor",
			"task_id", envelope.TaskID,
			"status", taskExecution.Status)
		return nil
	}

	task, err := e.buildTask(ctx, taskExecution, envelope)
	if err != nil {
		failErr := e.lifecycleManager.HandleTaskFailure(ctx, &taskRunnerDomainTask.Task{
			ID:                 envelope.TaskID,
			OrganizationID:     taskExecution.OrganizationID,
			NodeExecutionIDMap: make(map[string]uuid.UUID),
		}, err)
		if failErr != nil {
			slog.ErrorContext(ctx, "failed to mark task as failed after build error",
				"component", "task_executor",
				"task_id", envelope.TaskID,
				"error", failErr)
		}
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "task execution panicked",
				"component", "task_executor",
				"task_id", task.ID,
				"panic", r)
			failErr := e.lifecycleManager.HandleTaskFailure(ctx, task, ErrTaskExecutionPanic)
			if failErr != nil {
				slog.ErrorContext(ctx, "failed to handle task panic",
					"component", "task_executor",
					"task_id", task.ID,
					"error", failErr)
			}
		}
	}()

	maxConcurrentTasks := e.concurrencyLimitResolver.GetMaxConcurrentTasks(ctx, task.OrganizationID)
	semaphore, err := e.semaphoreManager.AcquireSemaphore(ctx, task.OrganizationID, task.ID, maxConcurrentTasks)
	if err != nil {
		failErr := e.lifecycleManager.HandleTaskFailure(ctx, task, err)
		if failErr != nil {
			slog.ErrorContext(ctx, "failed to handle task failure after semaphore error",
				"component", "task_executor",
				"task_id", task.ID,
				"error", failErr)
		}
		return err
	}

	semaphoreHeartbeatCtx, stopSemaphoreHeartbeat := context.WithCancel(ctx)
	defer stopSemaphoreHeartbeat()
	go e.runSemaphoreHeartbeat(semaphoreHeartbeatCtx, task.OrganizationID, task.ID)

	e.runTask(ctx, task, semaphore)
	claimReleased = e.releaseClaim(ctx, envelope.TaskID)
	return nil
}

func (e *TaskExecutor) runSemaphoreHeartbeat(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID) {
	runPeriodic(ctx, e.heartbeatInterval, func() (stop bool) {
		if err := e.semaphoreManager.HeartbeatSemaphore(ctx, organizationID, taskID); err != nil {
			slog.WarnContext(ctx, "failed to heartbeat semaphore",
				"component", "task_executor",
				"organization_id", organizationID,
				"task_id", taskID,
				"error", err)
		}

		return false
	})
}

func (e *TaskExecutor) buildTask(ctx context.Context, taskExecution *executionsDomainTaskExecutions.TaskExecution, envelope taskRunnerDomainTask.TaskEnvelope) (*taskRunnerDomainTask.Task, error) {
	if taskExecution.Nodes == nil || taskExecution.Edges == nil {
		return nil, ErrTaskDAGIsEmpty
	}

	dagInstance, err := dag.New(taskExecution.Nodes, taskExecution.Edges)
	if err != nil {
		return nil, err
	}

	nodeExecutions, err := e.nodeExecutionService.GetAllByTaskID(ctx, taskExecution.ID)
	nodeExecutionIDsByNodeID := make(map[string]uuid.UUID, len(nodeExecutions))
	if err == nil {
		for _, nodeExecution := range nodeExecutions {
			nodeExecutionIDsByNodeID[nodeExecution.NodeID] = nodeExecution.ID
		}
	}

	return &taskRunnerDomainTask.Task{
		ID:                  taskExecution.ID,
		OrganizationID:      taskExecution.OrganizationID,
		TriggeredByUserID:   taskExecution.TriggeredByUserID,
		ExecutionContext:    taskExecution.ExecutionContext,
		ContextItemID:       taskExecution.ContextItemID,
		DAG:                 dagInstance,
		StartTime:           taskExecution.StartedAt,
		NodeExecutionIDMap:  nodeExecutionIDsByNodeID,
		TriggerContext:      envelope.TriggerContext,
		PreviousNodeOutputs: envelope.PreviousNodeOutputs,
	}, nil
}

func (e *TaskExecutor) releaseClaim(ctx context.Context, taskID uuid.UUID) bool {
	if err := e.taskClaimService.ReleaseClaim(ctx, taskID, e.workerID); err != nil {
		slog.ErrorContext(ctx, "failed to release task claim",
			"component", "task_executor",
			"task_id", taskID,
			"worker_id", e.workerID,
			"error", err)
		return false
	}
	return true
}

func (e *TaskExecutor) runTask(ctx context.Context, task *taskRunnerDomainTask.Task, semaphore chan struct{}) {
	defer func() {
		e.contextManager.RemoveContext(task.ID)
		if err := e.semaphoreManager.ReleaseSemaphore(ctx, task.OrganizationID, task.ID, semaphore); err != nil {
			slog.ErrorContext(ctx, "failed to release semaphore",
				"component", "task_executor",
				"task_id", task.ID,
				"organization_id", task.OrganizationID,
				"error", err)
		}
	}()

	taskCtx := e.contextManager.CreateContext(task.ID, ctx)

	if e.maxExecutionTime > 0 {
		var cancelExecution context.CancelFunc
		taskCtx, cancelExecution = context.WithTimeout(taskCtx, e.maxExecutionTime)
		defer cancelExecution()
	}

	if err := e.lifecycleManager.StartTask(taskCtx, task); err != nil {
		if failErr := e.lifecycleManager.HandleTaskFailure(ctx, task, err); failErr != nil {
			slog.ErrorContext(ctx, "failed to handle task failure after start error",
				"component", "task_executor",
				"task_id", task.ID,
				"error", failErr)
		}
		return
	}

	if err := e.executionCoordinator.ExecuteTask(taskCtx, task); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = ErrTaskExecutionTimeLimit
		}
		if errors.Is(err, context.Canceled) {
			if cancelErr := e.lifecycleManager.HandleTaskCancellation(ctx, task); cancelErr != nil {
				slog.ErrorContext(ctx, "failed to handle task cancellation",
					"component", "task_executor",
					"task_id", task.ID,
					"error", cancelErr)
			}
			return
		}
		if failErr := e.lifecycleManager.HandleTaskFailure(ctx, task, err); failErr != nil {
			slog.ErrorContext(ctx, "failed to handle task failure",
				"component", "task_executor",
				"task_id", task.ID,
				"error", failErr)
		}
		return
	}

	if err := e.lifecycleManager.HandleTaskSuccess(ctx, task); err != nil {
		slog.ErrorContext(ctx, "failed to handle task success",
			"component", "task_executor",
			"task_id", task.ID,
			"error", err)
	}
}
