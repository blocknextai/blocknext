package taskrunner

import (
	"context"
	"log/slog"
	"strings"
	"time"

	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	executionsApplicationTaskClaims "github.com/blocknextai/platform-api/internal/executions/application/taskclaims"
	executionsApplicationTaskexecutions "github.com/blocknextai/platform-api/internal/executions/application/taskexecutions"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	taskRunnerApplicationContextresolver "github.com/blocknextai/platform-api/internal/taskrunner/application/contextresolver"
	taskRunnerDomain "github.com/blocknextai/platform-api/internal/taskrunner/domain"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/triggers/application/triggers"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
)

type TaskRecoveryService struct {
	taskExecutionService executionsApplicationTaskexecutions.TaskExecutionService
	taskClaimService     executionsApplicationTaskClaims.TaskClaimService
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService
	dispatcher           taskRunnerDomainTaskRunner.TaskDispatcher
	triggerService       triggersApplicationTriggers.TriggerService
	cronService          taskRunnerDomainTaskRunner.CronService
	taskService          taskRunnerDomainTaskRunner.TaskService
	contextResolver      taskRunnerApplicationContextresolver.ContextResolver
	staleClaimTimeout    time.Duration
	recoveryInterval     time.Duration
}

func NewTaskRecoveryService(
	taskExecutionService executionsApplicationTaskexecutions.TaskExecutionService,
	taskClaimService executionsApplicationTaskClaims.TaskClaimService,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	dispatcher taskRunnerDomainTaskRunner.TaskDispatcher,
	triggerService triggersApplicationTriggers.TriggerService,
	cronService taskRunnerDomainTaskRunner.CronService,
	taskService taskRunnerDomainTaskRunner.TaskService,
	contextResolver taskRunnerApplicationContextresolver.ContextResolver,
	staleClaimTimeout time.Duration,
	recoveryInterval time.Duration,
) *TaskRecoveryService {
	return &TaskRecoveryService{
		taskExecutionService: taskExecutionService,
		taskClaimService:     taskClaimService,
		nodeExecutionService: nodeExecutionService,
		dispatcher:           dispatcher,
		triggerService:       triggerService,
		cronService:          cronService,
		taskService:          taskService,
		contextResolver:      contextResolver,
		staleClaimTimeout:    staleClaimTimeout,
		recoveryInterval:     recoveryInterval,
	}
}

func (s *TaskRecoveryService) Start(ctx context.Context) {
	s.reconcileSchedules(ctx)
	s.runRecoveryPass(ctx)

	ticker := time.NewTicker(s.recoveryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.reconcileSchedules(ctx)
			s.runRecoveryPass(ctx)
		}
	}
}

func (s *TaskRecoveryService) runRecoveryPass(ctx context.Context) {
	releasedIDs, err := s.taskClaimService.ReleaseStaleClaims(ctx, s.staleClaimTimeout)
	if err != nil {
		slog.ErrorContext(ctx, "failed to release stale task claims",
			"component", "task_recovery_service",
			"error", err)
		return
	}

	if len(releasedIDs) == 0 {
		return
	}

	slog.InfoContext(ctx, "released stale task claims, recovering tasks",
		"component", "task_recovery_service",
		"count", len(releasedIDs))

	for _, taskID := range releasedIDs {
		taskExecution, err := s.taskExecutionService.GetByID(ctx, taskID)
		if err != nil {
			slog.ErrorContext(ctx, "failed to load task for recovery",
				"component", "task_recovery_service",
				"task_id", taskID,
				"error", err)
			continue
		}
		if !isActiveStatus(taskExecution.Status) {
			slog.InfoContext(ctx, "task not in active state, skipping recovery",
				"component", "task_recovery_service",
				"task_id", taskID,
				"status", taskExecution.Status)
			continue
		}
		s.recoverSingleTask(ctx, taskExecution)
	}
}

func isActiveStatus(status string) bool {
	return status == taskRunnerDomain.StatusPending.String() ||
		status == taskRunnerDomain.StatusRunning.String() ||
		status == taskRunnerDomain.StatusRetrying.String()
}

func (s *TaskRecoveryService) recoverSingleTask(ctx context.Context, taskExecution *taskexecutions.TaskExecution) {
	if taskExecution.Nodes == nil || taskExecution.Edges == nil {
		s.markTaskAsFailed(ctx, taskExecution, ErrTaskExecutionPanic.Error())
		return
	}

	if !s.hasIncompleteNodes(ctx, taskExecution) {
		s.determineTaskFinalState(ctx, taskExecution)
		return
	}

	envelope := taskRunnerDomainTask.TaskEnvelope{
		TaskID: taskExecution.ID,
	}
	if err := s.dispatcher.Dispatch(ctx, envelope); err != nil {
		slog.ErrorContext(ctx, "failed to re-dispatch task for recovery",
			"component", "task_recovery_service",
			"task_id", taskExecution.ID,
			"error", err)
		s.markTaskAsFailed(ctx, taskExecution, "recovery dispatch failed: "+err.Error())
	}
}

func (s *TaskRecoveryService) hasIncompleteNodes(ctx context.Context, taskExecution *taskexecutions.TaskExecution) bool {
	nodeExecutions, err := s.nodeExecutionService.GetAllByTaskID(ctx, taskExecution.ID)
	if err != nil {
		return true
	}

	if len(nodeExecutions) == 0 {
		return true
	}

	for _, nodeExecution := range nodeExecutions {
		if isActiveStatus(nodeExecution.Status) {
			return true
		}
	}

	return false
}

func (s *TaskRecoveryService) determineTaskFinalState(ctx context.Context, taskExecution *taskexecutions.TaskExecution) {
	nodeExecutions, err := s.nodeExecutionService.GetAllByTaskID(ctx, taskExecution.ID)
	if err != nil {
		s.markTaskAsFailed(ctx, taskExecution, ErrNoNodeExecutionsFound.Error())
		return
	}

	if len(nodeExecutions) == 0 {
		return
	}

	hasFailedNodes := false
	hasSuccessfulNodes := false

	for _, nodeExecution := range nodeExecutions {
		if nodeExecution.Status == taskRunnerDomain.StatusFailed.String() {
			hasFailedNodes = true
		} else if nodeExecution.Status == taskRunnerDomain.StatusSuccess.String() {
			hasSuccessfulNodes = true
		}
	}

	if hasFailedNodes {
		s.markTaskAsFailed(ctx, taskExecution, ErrSomeNodesFailed.Error())
	} else if hasSuccessfulNodes {
		s.markTaskAsSuccess(ctx, taskExecution)
	}
}

func (s *TaskRecoveryService) markTaskAsFailed(ctx context.Context, taskExecution *taskexecutions.TaskExecution, errorMessage string) {
	err := s.taskExecutionService.Update(
		ctx,
		taskExecution.ID,
		taskExecution.OrganizationID,
		taskRunnerDomain.StatusFailed.String(),
		&errorMessage,
		new(time.Now().UTC()),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to mark task as failed",
			"component", "task_recovery_service",
			"error", err,
			"task_id", taskExecution.ID,
			"organization_id", taskExecution.OrganizationID,
		)
	}
}

func (s *TaskRecoveryService) markTaskAsSuccess(ctx context.Context, taskExecution *taskexecutions.TaskExecution) {
	err := s.taskExecutionService.Update(
		ctx,
		taskExecution.ID,
		taskExecution.OrganizationID,
		taskRunnerDomain.StatusSuccess.String(),
		nil,
		new(time.Now().UTC()),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to mark task as success",
			"component", "task_recovery_service",
			"error", err,
			"task_id", taskExecution.ID,
			"organization_id", taskExecution.OrganizationID,
		)
	}
}

func (s *TaskRecoveryService) reconcileSchedules(ctx context.Context) {
	triggers, err := s.triggerService.GetAllActive(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get active triggers for schedule reconcile",
			"component", "task_recovery_service",
			"error", err)
		return
	}

	desired := make([]taskRunnerDomainTaskRunner.DesiredJob, 0, len(triggers))
	for _, trigger := range triggers {
		if trigger.Type != triggersDomainTriggers.TriggerTypeSchedule || trigger.CronPattern == nil {
			continue
		}
		if strings.TrimSpace(*trigger.CronPattern) == "" {
			continue
		}

		desired = append(desired, taskRunnerDomainTaskRunner.DesiredJob{
			TriggerID: trigger.ID,
			Pattern:   *trigger.CronPattern,
			Version:   trigger.UpdatedAt.UTC().Format(time.RFC3339Nano),
			Run:       s.buildScheduleJob(trigger),
		})
	}

	s.cronService.ReconcileJobs(desired)
}

func (s *TaskRecoveryService) buildScheduleJob(trigger *triggersDomainTriggers.Trigger) func() {
	return func() {
		bgCtx := context.Background()

		slog.InfoContext(bgCtx, "executing scheduled trigger",
			"component", "task_recovery_service",
			"trigger_id", trigger.ID)

		resolvedContext, err := s.contextResolver.ResolveContext(
			bgCtx, trigger.ExecutionContext, trigger.ContextItemID, trigger.OrganizationID,
		)
		if err != nil {
			slog.ErrorContext(bgCtx, "failed to resolve context for scheduled trigger",
				"component", "task_recovery_service",
				"trigger_id", trigger.ID,
				"error", err)
			return
		}

		ApplyRuntimeConfig(trigger.RuntimeConfig, resolvedContext.Nodes)

		var triggerContext *nodeEngineDomainAdapters.TriggerContext
		if trigger.RuntimeConfig != nil && strings.TrimSpace(trigger.RuntimeConfig.RuntimePrompt) != "" {
			sender := ""
			if trigger.TriggeredByUserID != nil {
				sender = trigger.TriggeredByUserID.String()
			}
			triggerContext = &nodeEngineDomainAdapters.TriggerContext{
				Source: taskRunnerDomainTask.TaskTriggerTypeSchedule.String(),
				Sender: sender,
				Prompt: trigger.RuntimeConfig.RuntimePrompt,
			}
		}

		_, err = s.taskService.ExecuteTask(
			bgCtx,
			trigger.TriggeredByUserID,
			trigger.OrganizationID,
			trigger.ExecutionContext,
			trigger.ContextItemID,
			taskRunnerDomainTask.TaskTriggerTypeSchedule,
			triggerContext,
			nil,
			nil,
			resolvedContext.Nodes,
			resolvedContext.Edges,
		)
		if err != nil {
			slog.ErrorContext(bgCtx, "failed to execute scheduled task",
				"component", "task_recovery_service",
				"trigger_id", trigger.ID,
				"error", err)
		}
	}
}
