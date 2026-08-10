package taskrunner

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/blocknextai/go-packages/dag"
	credentialOAuthApplicationRegenerate "github.com/blocknextai/platform-api/internal/credentialoauth/application/regenerate"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	executionsApplicationTaskexecutions "github.com/blocknextai/platform-api/internal/executions/application/taskexecutions"
	nodeEngineDomainExecutors "github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	taskRunnerDomain "github.com/blocknextai/platform-api/internal/taskrunner/domain"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/google/uuid"
)

type taskExecutionCoordinator struct {
	taskExecutionService   executionsApplicationTaskexecutions.TaskExecutionService
	nodeExecutionService   executionsApplicationNodeExecutions.NodeExecutionService
	tokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService
	nodeExecutor           *NodeExecutor
	eventPublisher         *EventPublisher
}

func NewTaskExecutionCoordinator(
	taskExecutionService executionsApplicationTaskexecutions.TaskExecutionService,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	tokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService,
	nodeExecutor *NodeExecutor,
	eventPublisher *EventPublisher,
) taskRunnerDomainTaskRunner.TaskExecutionCoordinator {
	return &taskExecutionCoordinator{
		taskExecutionService:   taskExecutionService,
		nodeExecutionService:   nodeExecutionService,
		tokenRegenerateService: tokenRegenerateService,
		nodeExecutor:           nodeExecutor,
		eventPublisher:         eventPublisher,
	}
}

func (c *taskExecutionCoordinator) ExecuteTask(ctx context.Context, task *taskRunnerDomainTask.Task) error {
	defer c.nodeExecutor.outputStore.Release(task.ID)
	return c.executeTaskNodes(ctx, task)
}

func (c *taskExecutionCoordinator) CancelExecution(ctx context.Context, taskID uuid.UUID) error {
	return nil
}

func (c *taskExecutionCoordinator) seedFromDatabase(ctx context.Context, task *taskRunnerDomainTask.Task) {
	taskID := task.ID
	nodeExecutions, err := c.nodeExecutionService.GetAllByTaskID(ctx, taskID)
	if err != nil {
		slog.WarnContext(ctx, "failed to load node executions for output restoration",
			"component", "task_execution_coordinator",
			"task_id", taskID,
			"error", err)
		return
	}

	restored := 0
	for _, nx := range nodeExecutions {
		task.NodeStatuses.Store(nx.NodeID, nx.Status)

		if nx.Status != taskRunnerDomain.StatusSuccess.String() {
			continue
		}
		if len(nx.Outputs) == 0 {
			continue
		}
		nodeKey := BuildNodeKey(nx.NodeType, nx.NodeID)
		c.nodeExecutor.outputStore.Store(taskID, nodeKey, nx.Outputs)
		restored++
	}

	if restored > 0 {
		slog.InfoContext(ctx, "restored node outputs from database",
			"component", "task_execution_coordinator",
			"task_id", taskID,
			"count", restored)
	}
}

func (c *taskExecutionCoordinator) executeTaskNodes(ctx context.Context, task *taskRunnerDomainTask.Task) error {
	if task.DAG == nil {
		return ErrTaskDAGIsEmpty
	}

	c.seedFromDatabase(ctx, task)

	startNodeIDs := task.DAG.StartNodes()
	if len(startNodeIDs) == 0 {
		return ErrNoStartNodesFound
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(startNodeIDs))
	ctxDone := ctx.Done()

	for _, nodeID := range startNodeIDs {
		select {
		case <-ctxDone:
			return context.Canceled
		default:
		}

		wg.Add(1)
		go func(nodeID string) {
			defer wg.Done()

			select {
			case <-ctxDone:
				errChan <- context.Canceled
				return
			default:
			}

			node := task.DAG.NodeByID(nodeID)
			if node == nil {
				return
			}
			if err := c.executeNode(ctx, task, node); err != nil {
				errChan <- err
			}
		}(nodeID)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *taskExecutionCoordinator) executeNode(ctx context.Context, task *taskRunnerDomainTask.Task, node *dag.Node) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if _, started := task.StartedNodes.LoadOrStore(node.ID, struct{}{}); started {
		return nil
	}

	nodeExecutionID, exists := task.NodeExecutionIDMap[node.ID]
	if !exists {
		if _, runnable := nodeEngineDomainExecutors.GetExecutor(node.NodeID); !runnable {
			return ErrExecutorNotFound
		}
		return ErrNodeExecutionNotFound
	}

	if task.PreviousNodeOutputs != nil && !routesItems(node) {
		nodeKey := BuildNodeKey(node.NodeID, node.ID)
		outputs, hasPreviousOutputs := task.PreviousNodeOutputs[nodeKey]
		if hasPreviousOutputs {
			c.nodeExecutor.outputStore.Store(task.ID, nodeKey, outputs)
			task.NodeStatuses.Store(node.ID, taskRunnerDomain.StatusSuccess.String())

			c.eventPublisher.PublishNodeEvent(
				ctx,
				task.ID,
				task.OrganizationID,
				task.ExecutionContext,
				task.ContextItemID,
				node.ID,
				node.NodeID,
				taskRunnerDomain.StatusSuccess,
				outputs,
				"",
				0,
			)

			slog.InfoContext(ctx, "Skipping node execution - using previous outputs",
				"component", "task_execution_coordinator",
				"task_id", task.ID,
				"node_id", node.ID,
				"node_execution_id", nodeExecutionID)

			return c.executeChildNodes(ctx, task, node)
		}
	}

	if node.Settings != nil && node.Settings.Disabled {
		c.markNodeAsSkipped(ctx, task, node, ErrNodeDisabled.Error())
		return c.executeChildNodes(ctx, task, node)
	}

	if err := c.nodeExecutor.ExecuteNode(ctx, task, node); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		if node.Settings != nil && node.Settings.ContinueOnError {
			return c.executeChildNodes(ctx, task, node)
		}
		return err
	}

	return c.executeChildNodes(ctx, task, node)
}

func (c *taskExecutionCoordinator) executeChildNodes(ctx context.Context, task *taskRunnerDomainTask.Task, node *dag.Node) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	childNodes := task.DAG.NodeChildren(node.ID)
	if len(childNodes) == 0 {
		return nil
	}

	readyChildren := make([]*dag.Node, 0)
	for _, childNodeID := range childNodes {
		childNode := task.DAG.NodeByID(childNodeID)
		if childNode == nil {
			continue
		}
		reached, err := branchReaches(c.nodeExecutor.outputStore, task.ID, task.DAG, node, childNodeID)
		if err != nil {
			return err
		}
		if !reached {
			c.skipUnreachable(ctx, task, childNode)
			continue
		}
		if !c.areParentNodesCompleted(task, childNode) {
			continue
		}
		readyChildren = append(readyChildren, childNode)
	}

	if len(readyChildren) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(readyChildren))

	for _, childNode := range readyChildren {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		wg.Add(1)
		go func(node *dag.Node) {
			defer wg.Done()
			if err := c.executeNode(ctx, task, node); err != nil {
				errChan <- err
			}
		}(childNode)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *taskExecutionCoordinator) markNodeAsSkipped(ctx context.Context, task *taskRunnerDomainTask.Task, node *dag.Node, reason string) {
	nodeExecutionID, exists := task.NodeExecutionIDMap[node.ID]
	if !exists {
		return
	}

	startTime := time.Now().UTC()

	if err := c.nodeExecutionService.Update(
		ctx,
		task.ID,
		nodeExecutionID,
		taskRunnerDomain.StatusRunning.String(),
		nil,
		nil,
		nil,
		nil,
		&startTime,
		nil,
	); err != nil {
		slog.ErrorContext(ctx, "Failed to update node state to running",
			"component", "task_execution_coordinator",
			"task_id", task.ID,
			"node_id", node.ID,
			"node_execution_id", nodeExecutionID,
			"error", err)
	}

	c.eventPublisher.PublishNodeEvent(
		ctx,
		task.ID,
		task.OrganizationID,
		task.ExecutionContext,
		task.ContextItemID,
		node.ID,
		node.NodeID,
		taskRunnerDomain.StatusPending,
		nil,
		"",
		0,
	)

	endTime := time.Now().UTC()
	executionTime := endTime.Sub(startTime).Milliseconds()
	task.NodeStatuses.Store(node.ID, taskRunnerDomain.StatusSkipped.String())

	var errorMessage *string
	if strings.TrimSpace(reason) != "" {
		errorMessage = new(reason)
	}

	if err := c.nodeExecutionService.Update(
		ctx,
		task.ID,
		nodeExecutionID,
		taskRunnerDomain.StatusSkipped.String(),
		nil,
		nil,
		nil,
		errorMessage,
		nil,
		&endTime,
	); err != nil {
		slog.ErrorContext(ctx, "Failed to update node state to skipped",
			"component", "task_execution_coordinator",
			"task_id", task.ID,
			"node_id", node.ID,
			"node_execution_id", nodeExecutionID,
			"error", err)
	}

	c.eventPublisher.PublishNodeEvent(
		ctx,
		task.ID,
		task.OrganizationID,
		task.ExecutionContext,
		task.ContextItemID,
		node.ID,
		node.NodeID,
		taskRunnerDomain.StatusSkipped,
		nil,
		reason,
		executionTime,
	)
}

func (c *taskExecutionCoordinator) areParentNodesCompleted(
	task *taskRunnerDomainTask.Task,
	node *dag.Node,
) bool {
	parentNodes := task.DAG.NodeParents(node.ID)
	if len(parentNodes) == 0 {
		return true
	}

	for _, parentNodeID := range parentNodes {
		parentNode := task.DAG.NodeByID(parentNodeID)
		if parentNode == nil {
			continue
		}

		status, exists := nodeStatus(task, parentNode.ID)
		if !exists {
			return false
		}

		isCompleted := status == taskRunnerDomain.StatusSuccess.String() ||
			status == taskRunnerDomain.StatusSkipped.String() ||
			(status == taskRunnerDomain.StatusFailed.String() && parentNode.Settings != nil && parentNode.Settings.ContinueOnError)

		if !isCompleted {
			return false
		}
	}

	return true
}

func (c *taskExecutionCoordinator) skipUnreachable(
	ctx context.Context,
	task *taskRunnerDomainTask.Task,
	root *dag.Node,
) {
	skipped := make(map[string]bool)

	var walk func(node *dag.Node)
	walk = func(node *dag.Node) {
		if skipped[node.ID] {
			return
		}
		skipped[node.ID] = true
		c.markNodeAsSkipped(ctx, task, node, "")

		for _, childID := range task.DAG.NodeChildren(node.ID) {
			child := task.DAG.NodeByID(childID)
			if child == nil {
				continue
			}
			if c.hasLiveParent(task, childID, skipped) {
				continue
			}
			walk(child)
		}
	}

	walk(root)
}

func (c *taskExecutionCoordinator) hasLiveParent(
	task *taskRunnerDomainTask.Task,
	nodeID string,
	skipped map[string]bool,
) bool {
	for _, parentID := range task.DAG.NodeParents(nodeID) {
		if skipped[parentID] {
			continue
		}
		if status, ok := nodeStatus(task, parentID); ok &&
			status == taskRunnerDomain.StatusSkipped.String() {
			continue
		}
		return true
	}
	return false
}

func nodeStatus(task *taskRunnerDomainTask.Task, nodeID string) (string, bool) {
	value, ok := task.NodeStatuses.Load(nodeID)
	if !ok {
		return "", false
	}
	status, ok := value.(string)
	return status, ok
}
