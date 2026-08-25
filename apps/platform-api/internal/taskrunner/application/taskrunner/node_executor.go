package taskrunner

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/dag"
	"github.com/blocknextai/go-packages/json"
	executionsApplicationNodeExecutions "github.com/blocknextai/platform-api/internal/executions/application/nodeexecutions"
	functionCallingPkg "github.com/blocknextai/platform-api/internal/llm/functioncalling"
	nodeEngineDomainExecutors "github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	taskRunnerDomain "github.com/blocknextai/platform-api/internal/taskrunner/domain"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
)

type NodeExecutor struct {
	functionCallingService functionCallingPkg.FunctionCallingService
	dataProcessor          *DataProcessor
	credentialProcessor    *CredentialProcessor
	eventPublisher         *EventPublisher
	nodeExecutionService   executionsApplicationNodeExecutions.NodeExecutionService
	outputStore            *OutputStore
}

func NewNodeExecutor(
	functionCallingService functionCallingPkg.FunctionCallingService,
	dataProcessor *DataProcessor,
	credentialProcessor *CredentialProcessor,
	eventPublisher *EventPublisher,
	nodeExecutionService executionsApplicationNodeExecutions.NodeExecutionService,
	outputStore *OutputStore,
) *NodeExecutor {
	return &NodeExecutor{
		functionCallingService: functionCallingService,
		dataProcessor:          dataProcessor,
		credentialProcessor:    credentialProcessor,
		eventPublisher:         eventPublisher,
		nodeExecutionService:   nodeExecutionService,
		outputStore:            outputStore,
	}
}

func (e *NodeExecutor) ExecuteNode(ctx context.Context, task *taskRunnerDomainTask.Task, node *dag.Node) error {
	startTime := time.Now().UTC()

	if _, exists := task.NodeExecutionIDMap[node.ID]; !exists {
		return ErrNodeExecutionIDNotFound
	}

	e.publishNodeEvent(ctx, task, node, taskRunnerDomain.StatusPending, nil, "", 0)
	e.publishNodeEvent(ctx, task, node, taskRunnerDomain.StatusRunning, nil, "", 0)

	executor, ok := nodeEngineDomainExecutors.GetExecutor(node.NodeID)
	if !ok {
		return ErrExecutorNotFound
	}

	maxRetries, retryDelay, timeout := e.extractSettings(node)

	var processedCredentials map[string]any
	if len(node.Credentials) > 0 {
		var credErr error
		processedCredentials, credErr = e.credentialProcessor.ProcessCredentials(
			ctx, task.OrganizationID, node.NodeID, node.Credentials,
		)
		if credErr != nil {
			e.finalizeNode(ctx, task, node, nil, nil, credErr, startTime, taskRunnerDomain.StatusFailed)
			return credErr
		}
	}

	view := e.buildBranchView(task, node)
	processedDataList := e.prepareAndProcessNodeData(task, node, view)

	e.updateNodeExecution(ctx, task.ID, task.NodeExecutionIDMap[node.ID], taskRunnerDomain.StatusRunning,
		processedDataList, nil, nil, nil, &startTime, nil)

	return e.executeWithRetries(ctx, task, node, executor, processedCredentials,
		processedDataList, node.Parameters, maxRetries, retryDelay, timeout, startTime, view)
}

func (e *NodeExecutor) executeWithRetries(
	ctx context.Context,
	task *taskRunnerDomainTask.Task,
	node *dag.Node,
	executor nodeEngineDomainExecutors.ExecutorManager,
	credentials map[string]any,
	processedDataList []map[string]any,
	parameters map[string]any,
	maxRetries, retryDelay, timeout int,
	startTime time.Time,
	view *BranchView,
) error {
	var lastErr error
	var lastFunctionCallingOutputs []map[string]any

	for attempt := 0; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			e.finalizeNode(ctx, task, node, nil, lastFunctionCallingOutputs, ctx.Err(), startTime, taskRunnerDomain.StatusCancelled)
			return ctx.Err()
		default:
			outputs, functionCallingOutputs, branches, err := e.executeOperation(ctx, task, executor, credentials, processedDataList, parameters, timeout, view)
			lastFunctionCallingOutputs = functionCallingOutputs

			if err == nil {
				e.storeNodeOutputs(task.ID, node, outputs, processedDataList, branches, view)
				e.finalizeNode(ctx, task, node, outputs, functionCallingOutputs, nil, startTime, taskRunnerDomain.StatusSuccess)
				return nil
			}

			lastErr = err
			if errors.Is(lastErr, context.Canceled) {
				e.finalizeNode(ctx, task, node, nil, functionCallingOutputs, lastErr, startTime, taskRunnerDomain.StatusCancelled)
				return lastErr
			}

			if attempt < maxRetries {
				e.publishNodeEvent(ctx, task, node, taskRunnerDomain.StatusRetrying, nil, err.Error(),
					time.Since(startTime).Milliseconds())

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Duration(retryDelay) * time.Millisecond):
					continue
				}
			}
		}
	}

	e.finalizeNode(ctx, task, node, nil, lastFunctionCallingOutputs, lastErr, startTime, taskRunnerDomain.StatusFailed)
	return lastErr
}

func (e *NodeExecutor) executeOperation(
	ctx context.Context,
	task *taskRunnerDomainTask.Task,
	executor nodeEngineDomainExecutors.ExecutorManager,
	credentials map[string]any,
	processedDataList []map[string]any,
	parameters map[string]any,
	timeout int,
	view *BranchView,
) ([]map[string]any, []map[string]any, map[string][]int, error) {
	attemptCtx := ctx
	if timeout > 0 {
		var cancel context.CancelFunc
		attemptCtx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Millisecond)
		defer cancel()
	}

	if processedDataList == nil {
		processedDataList = []map[string]any{{}}
	}

	executorID := executor.GetID()
	var functionCallingOutputs []map[string]any
	fc, ok := functioncalling.GetFunctionCalling(executorID)
	if ok && e.functionCallingService != nil {
		nodeFunctionCallingOutputs, fcErr := e.functionCallingService.ExecuteWithContext(attemptCtx, processedDataList, fc.GetFunctionDeclarations())
		if fcErr != nil {
			slog.ErrorContext(ctx, "Function calling failed",
				"component", "node_executor",
				"executor_id", executorID,
				"error", fcErr)
		}
		functionCallingOutputs = append(functionCallingOutputs, nodeFunctionCallingOutputs...)

		if len(functionCallingOutputs) > 0 {
			processedDataList = functionCallingOutputs
		}
	}

	if len(parameters) > 0 {
		var triggerData map[string]any
		if task != nil && task.TriggerContext != nil {
			if err := json.ArgsToStruct(task.TriggerContext, &triggerData); err != nil {
				triggerData = nil
			}
		}

		for i, data := range processedDataList {
			resolvedParameters := e.dataProcessor.ProcessNodeData(parameters, i, triggerData, view)
			for k, v := range resolvedParameters {
				if v == nil {
					continue
				}
				if s, ok := v.(string); ok && s == "" {
					continue
				}
				data[k] = v
			}
			processedDataList[i] = data
		}
	}

	if branching, ok := executor.(nodeEngineDomainExecutors.BranchingExecutor); ok {
		executorOutputs, branches, err := branching.ExecuteBranches(attemptCtx, credentials, processedDataList)
		return executorOutputs, functionCallingOutputs, branches, err
	}

	executorOutputs, err := executor.ExecuteWithContext(attemptCtx, credentials, processedDataList)
	return executorOutputs, functionCallingOutputs, nil, err
}

func (e *NodeExecutor) prepareAndProcessNodeData(task *taskRunnerDomainTask.Task, node *dag.Node, view *BranchView) []map[string]any {
	nodeDataMap := make(map[string]any)

	if node.Instruction != "" {
		nodeDataMap["instruction"] = node.Instruction
	}
	if node.RuntimeInstruction != "" {
		nodeDataMap["runtimeInstruction"] = node.RuntimeInstruction
	}
	if node.RuntimePrompt != "" {
		nodeDataMap["runtimePrompt"] = node.RuntimePrompt
	}

	var triggerData map[string]any
	if task.TriggerContext != nil {
		if err := json.ArgsToStruct(task.TriggerContext, &triggerData); err != nil {
			triggerData = nil
		}
	}

	parentNodes := task.DAG.NodeParents(node.ID)
	maxItems := e.calculateMaxItems(task, parentNodes, view)
	executionDataList := e.createExecutionDataList(nodeDataMap, maxItems)

	for i, data := range executionDataList {
		executionDataList[i] = e.dataProcessor.ProcessNodeData(data, i, triggerData, view)
	}

	return executionDataList
}

func (e *NodeExecutor) calculateMaxItems(task *taskRunnerDomainTask.Task, parentNodes []string, view *BranchView) int {
	maxItems := 1
	var builder strings.Builder
	for _, parentID := range parentNodes {
		parentNode := task.DAG.NodeByID(parentID)
		if parentNode == nil {
			continue
		}
		nodeKey := BuildNodeKeyWithBuilder(&builder, parentNode.NodeID, parentNode.ID)
		if outputs, ok := e.outputStore.Get(view.Task(), view.StoreKey(nodeKey)); ok && len(outputs) > maxItems {
			maxItems = len(outputs)
		}
	}
	return maxItems
}

func (e *NodeExecutor) createExecutionDataList(nodeData map[string]any, maxItems int) []map[string]any {
	executionDataList := make([]map[string]any, 0, maxItems)
	for range maxItems {
		executionData := make(map[string]any)
		maps.Copy(executionData, nodeData)
		executionDataList = append(executionDataList, executionData)
	}
	return executionDataList
}

func (e *NodeExecutor) finalizeNode(
	ctx context.Context,
	task *taskRunnerDomainTask.Task,
	node *dag.Node,
	outputs []map[string]any,
	functionCallingOutputs []map[string]any,
	err error,
	startTime time.Time,
	status taskRunnerDomain.Status,
) {
	executionTime := time.Since(startTime).Milliseconds()
	task.NodeStatuses.Store(node.ID, status.String())
	nodeExecutionID, exists := task.NodeExecutionIDMap[node.ID]
	if !exists {
		return
	}

	var errorMsg *string
	if err != nil {
		msg := err.Error()
		errorMsg = new(msg)
	}

	e.updateNodeExecution(ctx, task.ID, nodeExecutionID, status,
		nil, outputs, functionCallingOutputs,
		errorMsg, nil, new(time.Now().UTC()))

	errorStr := ""
	if errorMsg != nil {
		errorStr = *errorMsg
	}
	e.publishNodeEvent(ctx, task, node, status, outputs, errorStr, executionTime)
}

func (e *NodeExecutor) storeNodeOutputs(
	taskID uuid.UUID,
	node *dag.Node,
	outputs []map[string]any,
	items []map[string]any,
	branches map[string][]int,
	view *BranchView,
) {
	nodeKey := BuildNodeKey(node.NodeID, node.ID)
	e.outputStore.Store(taskID, nodeKey, outputs)

	for handle, indexes := range branches {
		branchItems := make([]map[string]any, 0, len(indexes))
		absolute := make([]int, 0, len(indexes))
		for _, index := range indexes {
			if index < 0 || index >= len(items) {
				continue
			}
			branchItems = append(branchItems, items[index])
			absolute = append(absolute, view.AbsoluteIndex(index))
		}
		e.outputStore.StoreBranch(taskID, nodeKey, handle, branchItems, absolute)
	}
}

func (e *NodeExecutor) extractSettings(node *dag.Node) (maxRetries, retryDelay, timeout int) {
	if node.Settings == nil {
		return 0, 0, 0
	}
	if node.Settings.MaxRetries > 0 {
		maxRetries = int(node.Settings.MaxRetries)
	}
	if node.Settings.RetryDelay > 0 {
		retryDelay = int(node.Settings.RetryDelay)
	}
	if node.Settings.Timeout > 0 {
		timeout = int(node.Settings.Timeout)
	}
	return
}

func (e *NodeExecutor) updateNodeExecution(
	ctx context.Context,
	taskID uuid.UUID,
	nodeExecutionID uuid.UUID,
	status taskRunnerDomain.Status,
	inputs []map[string]any,
	outputs []map[string]any,
	functionCallingOutputs []map[string]any,
	errorMsg *string,
	startTime *time.Time,
	endTime *time.Time,
) {
	if err := e.nodeExecutionService.Update(
		ctx, taskID, nodeExecutionID, status.String(),
		inputs, outputs, functionCallingOutputs, errorMsg, startTime, endTime,
	); err != nil {
		slog.ErrorContext(ctx, "Failed to update node execution",
			"component", "node_executor",
			"task_id", taskID,
			"node_execution_id", nodeExecutionID,
			"status", status.String(),
			"error", err)
	}
}

func (e *NodeExecutor) publishNodeEvent(
	ctx context.Context,
	task *taskRunnerDomainTask.Task,
	node *dag.Node,
	status taskRunnerDomain.Status,
	output []map[string]any,
	errorMsg string,
	duration int64,
) {
	e.eventPublisher.PublishNodeEvent(
		ctx, task.ID, task.OrganizationID, task.ExecutionContext, task.ContextItemID,
		node.ID, node.NodeID, status, output, errorMsg, duration,
	)
}
