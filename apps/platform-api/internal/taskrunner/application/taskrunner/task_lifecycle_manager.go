package taskrunner

import (
	"context"
	"log/slog"
	"time"

	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	executionsApplicationTaskexecutions "github.com/blocknextai/platform-api/internal/executions/application/taskexecutions"
	nodeEngineDomainExecutors "github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	taskRunnerDomain "github.com/blocknextai/platform-api/internal/taskrunner/domain"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/google/uuid"
)

type taskLifecycleManager struct {
	taskExecutionService executionsApplicationTaskexecutions.TaskExecutionService
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService
	eventPublisher       *EventPublisher
}

func NewTaskLifecycleManager(
	taskExecutionService executionsApplicationTaskexecutions.TaskExecutionService,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	eventPublisher *EventPublisher,
) taskRunnerDomainTaskRunner.TaskLifecycleManager {
	return &taskLifecycleManager{
		taskExecutionService: taskExecutionService,
		nodeExecutionService: nodeExecutionService,
		eventPublisher:       eventPublisher,
	}
}

func (m *taskLifecycleManager) StartTask(ctx context.Context, task *taskRunnerDomainTask.Task) error {
	startTime := time.Now().UTC()
	task.StartTime = new(startTime)
	executionTime := startTime.Sub(*task.StartTime).Milliseconds()

	if err := m.taskExecutionService.Update(
		ctx,
		task.ID,
		task.OrganizationID,
		taskRunnerDomain.StatusRunning.String(),
		nil,
		nil,
	); err != nil {
		return err
	}

	if err := m.CreateNodeExecutions(ctx, task); err != nil {
		return err
	}

	m.eventPublisher.PublishTaskEvent(
		ctx,
		task.ID,
		task.OrganizationID,
		task.ExecutionContext,
		task.ContextItemID,
		taskRunnerDomain.StatusRunning,
		"",
		executionTime,
	)

	return nil
}

func (m *taskLifecycleManager) HandleTaskSuccess(ctx context.Context, task *taskRunnerDomainTask.Task) error {
	endTime := time.Now().UTC()
	task.EndTime = new(endTime)
	executionTime := endTime.Sub(*task.StartTime).Milliseconds()

	if err := m.taskExecutionService.Update(
		ctx,
		task.ID,
		task.OrganizationID,
		taskRunnerDomain.StatusSuccess.String(),
		nil,
		&endTime,
	); err != nil {
		return err
	}

	m.eventPublisher.PublishTaskEvent(
		ctx,
		task.ID,
		task.OrganizationID,
		task.ExecutionContext,
		task.ContextItemID,
		taskRunnerDomain.StatusSuccess,
		"",
		executionTime,
	)

	return nil
}

func (m *taskLifecycleManager) HandleTaskFailure(ctx context.Context, task *taskRunnerDomainTask.Task, err error) error {
	endTime := time.Now().UTC()
	task.EndTime = new(endTime)
	executionTime := endTime.Sub(*task.StartTime).Milliseconds()

	if cleanupErr := m.failRunningNodes(ctx, task, err.Error()); cleanupErr != nil {
		slog.ErrorContext(ctx, "Failed to fail running nodes",
			"component", "task_lifecycle_manager",
			"task_id", task.ID,
			"organization_id", task.OrganizationID,
			"error", cleanupErr)
	}

	errorMessage := err.Error()
	if err := m.taskExecutionService.Update(
		ctx,
		task.ID,
		task.OrganizationID,
		taskRunnerDomain.StatusFailed.String(),
		&errorMessage,
		&endTime,
	); err != nil {
		return err
	}

	m.eventPublisher.PublishTaskEvent(
		ctx,
		task.ID,
		task.OrganizationID,
		task.ExecutionContext,
		task.ContextItemID,
		taskRunnerDomain.StatusFailed,
		errorMessage,
		executionTime,
	)

	return nil
}

func (m *taskLifecycleManager) HandleTaskCancellation(ctx context.Context, task *taskRunnerDomainTask.Task) error {
	endTime := time.Now().UTC()
	task.EndTime = new(endTime)
	executionTime := endTime.Sub(*task.StartTime).Milliseconds()

	if err := m.taskExecutionService.Update(
		ctx,
		task.ID,
		task.OrganizationID,
		taskRunnerDomain.StatusCancelled.String(),
		nil,
		&endTime,
	); err != nil {
		return err
	}

	if err := m.cancelRunningNodes(ctx, task); err != nil {
		return err
	}

	m.eventPublisher.PublishTaskEvent(
		ctx,
		task.ID,
		task.OrganizationID,
		task.ExecutionContext,
		task.ContextItemID,
		taskRunnerDomain.StatusCancelled,
		"",
		executionTime,
	)

	return nil
}

func (m *taskLifecycleManager) CreateNodeExecutions(ctx context.Context, task *taskRunnerDomainTask.Task) error {
	if task.DAG == nil {
		return ErrTaskDAGIsEmpty
	}

	allNodes := task.DAG.Nodes()
	if len(allNodes) == 0 {
		return ErrNoNodesFoundInWorkflow
	}

	if task.NodeExecutionIDMap == nil {
		task.NodeExecutionIDMap = make(map[string]uuid.UUID)
	}

	utcNow := time.Now().UTC()
	for _, node := range allNodes {
		if _, exists := task.NodeExecutionIDMap[node.ID]; exists {
			continue
		}

		if _, runnable := nodeEngineDomainExecutors.GetExecutor(node.NodeID); !runnable {
			continue
		}

		nodeKey := BuildNodeKey(node.NodeID, node.ID)

		if outputs, hasPreviousOutputs := task.PreviousNodeOutputs[nodeKey]; hasPreviousOutputs {
			nodeExecutionID, err := m.nodeExecutionService.Create(
				ctx,
				task.ID,
				node.NodeID,
				node.ID,
				taskRunnerDomain.StatusSuccess.String(),
				nil,
				nil,
				&utcNow,
			)
			if err != nil {
				return err
			}
			task.NodeExecutionIDMap[node.ID] = nodeExecutionID

			if err := m.nodeExecutionService.Update(
				ctx,
				task.ID,
				nodeExecutionID,
				taskRunnerDomain.StatusSuccess.String(),
				nil,
				outputs,
				nil,
				nil,
				nil,
				&utcNow,
			); err != nil {
				slog.ErrorContext(ctx, "Failed to update node execution with previous outputs",
					"component", "task_lifecycle_manager",
					"task_id", task.ID,
					"node_id", node.ID,
					"node_execution_id", nodeExecutionID,
					"error", err)
			}
		} else {
			nodeExecutionID, err := m.nodeExecutionService.Create(
				ctx,
				task.ID,
				node.NodeID,
				node.ID,
				taskRunnerDomain.StatusPending.String(),
				nil,
				nil,
				nil,
			)
			if err != nil {
				return err
			}
			task.NodeExecutionIDMap[node.ID] = nodeExecutionID
		}
	}

	return nil
}

func (m *taskLifecycleManager) cancelRunningNodes(ctx context.Context, task *taskRunnerDomainTask.Task) error {
	return m.terminateRunningNodes(ctx, task, taskRunnerDomain.StatusCancelled, nil)
}

func (m *taskLifecycleManager) failRunningNodes(ctx context.Context, task *taskRunnerDomainTask.Task, errorMessage string) error {
	return m.terminateRunningNodes(ctx, task, taskRunnerDomain.StatusFailed, &errorMessage)
}

func (m *taskLifecycleManager) terminateRunningNodes(
	ctx context.Context,
	task *taskRunnerDomainTask.Task,
	status taskRunnerDomain.Status,
	errorMessage *string,
) error {
	nodeExecutions, err := m.nodeExecutionService.GetAllByTaskID(ctx, task.ID)
	if err != nil {
		return err
	}

	utcNow := time.Now().UTC()
	for _, nodeExecution := range nodeExecutions {
		if nodeExecution.Status != taskRunnerDomain.StatusRunning.String() &&
			nodeExecution.Status != taskRunnerDomain.StatusPending.String() {
			continue
		}

		var inputs []map[string]any
		var outputs []map[string]any
		var functionCallingOutputs []map[string]any
		if errorMessage != nil {
			inputs = nodeExecution.Inputs
			outputs = nodeExecution.Outputs
			functionCallingOutputs = nodeExecution.FunctionCallingOutputs
		}

		if err := m.nodeExecutionService.Update(
			ctx,
			task.ID,
			nodeExecution.ID,
			status.String(),
			inputs,
			outputs,
			functionCallingOutputs,
			errorMessage,
			nil,
			&utcNow,
		); err != nil {
			continue
		}
	}

	return nil
}
