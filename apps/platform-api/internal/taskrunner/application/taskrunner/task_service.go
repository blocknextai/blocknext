package taskrunner

import (
	"context"
	"log/slog"
	"time"

	"github.com/blocknextai/go-packages/dag"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	executionsApplicationTaskexecutions "github.com/blocknextai/platform-api/internal/executions/application/taskexecutions"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	taskRunnerApplicationContextresolver "github.com/blocknextai/platform-api/internal/taskrunner/application/contextresolver"
	taskRunnerDomain "github.com/blocknextai/platform-api/internal/taskrunner/domain"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/triggers/application/triggers"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

type TaskService struct {
	contextResolver      taskRunnerApplicationContextresolver.ContextResolver
	executionCoordinator taskRunnerDomainTaskRunner.TaskExecutionCoordinator
	dispatcher           taskRunnerDomainTaskRunner.TaskDispatcher
	contextManager       taskRunnerDomainTaskRunner.ContextManager
	taskExecutionService executionsApplicationTaskexecutions.TaskExecutionService
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService
	triggerService       triggersApplicationTriggers.TriggerService
}

func NewTaskService(
	contextResolver taskRunnerApplicationContextresolver.ContextResolver,
	executionCoordinator taskRunnerDomainTaskRunner.TaskExecutionCoordinator,
	dispatcher taskRunnerDomainTaskRunner.TaskDispatcher,
	contextManager taskRunnerDomainTaskRunner.ContextManager,
	taskExecutionService executionsApplicationTaskexecutions.TaskExecutionService,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	triggerService triggersApplicationTriggers.TriggerService,
) taskRunnerDomainTaskRunner.TaskService {
	return &TaskService{
		contextResolver:      contextResolver,
		executionCoordinator: executionCoordinator,
		dispatcher:           dispatcher,
		contextManager:       contextManager,
		taskExecutionService: taskExecutionService,
		nodeExecutionService: nodeExecutionService,
		triggerService:       triggerService,
	}
}

func (s *TaskService) ExecuteTask(
	ctx context.Context,
	triggeredByUserID *uuid.UUID,
	organizationID uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	triggerType taskRunnerDomainTask.TaskTriggerType,
	triggerContext *nodeEngineDomainAdapters.TriggerContext,
	cronPattern *string,
	runtimeConfig *triggersDomainTriggers.RuntimeConfig,
	nodes []dag.Node,
	edges []dag.Edge,
) (*taskRunnerDomainTask.Task, error) {
	if triggerType == taskRunnerDomainTask.TaskTriggerTypeSchedule && cronPattern != nil {
		return s.registerScheduleTrigger(ctx, triggeredByUserID, organizationID, executionContext, contextItemID, cronPattern, runtimeConfig)
	}

	if triggerType == taskRunnerDomainTask.TaskTriggerTypeWebhook && triggerContext == nil {
		return s.registerWebhookTrigger(ctx, triggeredByUserID, organizationID, executionContext, contextItemID, runtimeConfig)
	}

	return s.submitTask(ctx, triggeredByUserID, organizationID, executionContext, contextItemID, triggerType, triggerContext, nodes, edges, nil)
}

func (s *TaskService) CancelTask(ctx context.Context, triggeredByUserID *uuid.UUID, organizationID uuid.UUID, taskID uuid.UUID) error {
	taskExecution, err := s.taskExecutionService.GetByIDAndOrganizationID(ctx, taskID, organizationID)
	if err != nil {
		return ErrTaskNotFound
	}

	if taskExecution.Status == taskRunnerDomain.StatusSuccess.String() ||
		taskExecution.Status == taskRunnerDomain.StatusFailed.String() ||
		taskExecution.Status == taskRunnerDomain.StatusCancelled.String() {
		return ErrTaskAlreadyCompleted
	}

	s.contextManager.CancelContext(taskID)

	nodeExecutions, err := s.nodeExecutionService.GetAllByTaskID(ctx, taskID)
	if err == nil {
		taskCancelledMessage := ErrTaskCancelled.Error()
		completedAt := time.Now().UTC()
		for _, nodeExecution := range nodeExecutions {
			if nodeExecution.Status == taskRunnerDomain.StatusRunning.String() ||
				nodeExecution.Status == taskRunnerDomain.StatusPending.String() {
				err := s.nodeExecutionService.Update(
					ctx,
					taskID,
					nodeExecution.ID,
					taskRunnerDomain.StatusCancelled.String(),
					nodeExecution.Inputs,
					nodeExecution.Outputs,
					nodeExecution.FunctionCallingOutputs,
					&taskCancelledMessage,
					nil,
					&completedAt,
				)
				if err != nil {
					slog.ErrorContext(ctx, "Failed to update node execution",
						"component", "task_service",
						"node_execution_id", nodeExecution.ID,
						"error", err)
				}
			}
		}
	}

	if err := s.taskExecutionService.Update(ctx, taskID, organizationID, taskRunnerDomain.StatusCancelled.String(), nil, new(time.Now().UTC())); err != nil {
		return err
	}

	if err := s.executionCoordinator.CancelExecution(ctx, taskID); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) RerunAll(
	ctx context.Context,
	triggeredByUserID *uuid.UUID,
	organizationID uuid.UUID,
	taskID uuid.UUID,
	triggerType taskRunnerDomainTask.TaskTriggerType,
) (*taskRunnerDomainTask.Task, error) {
	originalTaskExecution, err := s.taskExecutionService.GetByIDAndOrganizationID(ctx, taskID, organizationID)
	if err != nil {
		return nil, err
	}

	return s.submitTask(ctx, triggeredByUserID, organizationID, originalTaskExecution.ExecutionContext, originalTaskExecution.ContextItemID, triggerType, nil, originalTaskExecution.Nodes, originalTaskExecution.Edges, nil)
}

func (s *TaskService) RerunFailed(
	ctx context.Context,
	triggeredByUserID *uuid.UUID,
	organizationID uuid.UUID,
	taskID uuid.UUID,
	triggerType taskRunnerDomainTask.TaskTriggerType,
) (*taskRunnerDomainTask.Task, error) {
	originalTaskExecution, err := s.taskExecutionService.GetByIDAndOrganizationID(ctx, taskID, organizationID)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	if originalTaskExecution.Status != taskRunnerDomain.StatusFailed.String() {
		return nil, ErrTaskNotFailed
	}

	nodeExecutions, err := s.nodeExecutionService.GetAllByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	successNodeOutputs := make(map[string][]map[string]any)
	for _, ne := range nodeExecutions {
		if ne.Status == taskRunnerDomain.StatusSuccess.String() && ne.Outputs != nil {
			nodeKey := BuildNodeKey(ne.NodeType, ne.NodeID)
			successNodeOutputs[nodeKey] = ne.Outputs
		}
	}

	return s.submitTask(ctx, triggeredByUserID, organizationID, originalTaskExecution.ExecutionContext, originalTaskExecution.ContextItemID, triggerType, nil, originalTaskExecution.Nodes, originalTaskExecution.Edges, successNodeOutputs)
}

func (s *TaskService) registerScheduleTrigger(
	ctx context.Context,
	triggeredByUserID *uuid.UUID,
	organizationID uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	cronPattern *string,
	runtimeConfig *triggersDomainTriggers.RuntimeConfig,
) (*taskRunnerDomainTask.Task, error) {
	defaultTz := "UTC"
	timezone := new(defaultTz)

	if _, _, err := s.triggerService.Create(
		ctx, organizationID, triggeredByUserID, executionContext, contextItemID,
		triggersDomainTriggers.TriggerTypeSchedule, cronPattern, timezone, runtimeConfig,
	); err != nil {
		return nil, err
	}

	return &taskRunnerDomainTask.Task{
		OrganizationID:     organizationID,
		TriggeredByUserID:  triggeredByUserID,
		ExecutionContext:   executionContext,
		ContextItemID:      contextItemID,
		Status:             taskRunnerDomain.StatusScheduled,
		NodeExecutionIDMap: make(map[string]uuid.UUID),
	}, nil
}

func (s *TaskService) registerWebhookTrigger(
	ctx context.Context,
	triggeredByUserID *uuid.UUID,
	organizationID uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	runtimeConfig *triggersDomainTriggers.RuntimeConfig,
) (*taskRunnerDomainTask.Task, error) {
	_, plainToken, err := s.triggerService.Create(
		ctx, organizationID, triggeredByUserID, executionContext, contextItemID,
		triggersDomainTriggers.TriggerTypeWebhook, nil, nil, runtimeConfig,
	)
	if err != nil {
		return nil, err
	}

	return &taskRunnerDomainTask.Task{
		OrganizationID:     organizationID,
		TriggeredByUserID:  triggeredByUserID,
		ExecutionContext:   executionContext,
		ContextItemID:      contextItemID,
		Status:             taskRunnerDomain.StatusScheduled,
		WebhookToken:       new(plainToken),
		NodeExecutionIDMap: make(map[string]uuid.UUID),
	}, nil
}

func (s *TaskService) submitTask(
	ctx context.Context,
	triggeredByUserID *uuid.UUID,
	organizationID uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	triggerType taskRunnerDomainTask.TaskTriggerType,
	triggerContext *nodeEngineDomainAdapters.TriggerContext,
	nodes []dag.Node,
	edges []dag.Edge,
	previousNodeOutputs map[string][]map[string]any,
) (*taskRunnerDomainTask.Task, error) {
	dagInstance, err := dag.New(nodes, edges)
	if err != nil {
		return nil, err
	}

	taskID := bnuuid.NewV7()
	startTime := time.Now().UTC()

	err = s.taskExecutionService.Create(
		ctx, taskID, organizationID, resolveTriggeredByUserID(triggerType, triggeredByUserID), nil,
		executionContext, contextItemID,
		taskRunnerDomain.StatusPending.String(), mapTriggerTypeToExecutionType(triggerType),
		nil, nodes, edges, &startTime,
	)
	if err != nil {
		return nil, err
	}

	task := &taskRunnerDomainTask.Task{
		ID:                  taskID,
		TriggeredByUserID:   triggeredByUserID,
		OrganizationID:      organizationID,
		ExecutionContext:    executionContext,
		ContextItemID:       contextItemID,
		DAG:                 dagInstance,
		StartTime:           new(startTime),
		NodeExecutionIDMap:  make(map[string]uuid.UUID),
		PreviousNodeOutputs: previousNodeOutputs,
		TriggerContext:      triggerContext,
	}

	envelope := taskRunnerDomainTask.TaskEnvelope{
		TaskID:              taskID,
		TriggerContext:      triggerContext,
		PreviousNodeOutputs: previousNodeOutputs,
	}

	if err := s.dispatcher.Dispatch(ctx, envelope); err != nil {
		return nil, err
	}

	return task, nil
}
