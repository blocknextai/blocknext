package taskrunner

import (
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/task"
	triggersContract "github.com/blocknextai/platform-api/internal/modules/triggers/contract"
)

func mapTriggerTypeToExecutionType(triggerType taskRunnerDomainTask.TaskTriggerType) executionsContract.ExecutionType {
	switch triggerType {
	case taskRunnerDomainTask.TaskTriggerTypeSchedule:
		return executionsContract.ExecutionTypeSchedule
	case taskRunnerDomainTask.TaskTriggerTypeWebhook:
		return executionsContract.ExecutionTypeWebhook
	case taskRunnerDomainTask.TaskTriggerTypeAPI:
		return executionsContract.ExecutionTypeAPI
	default:
		return executionsContract.ExecutionTypeManual
	}
}

func resolveTriggeredByUserID(triggerType taskRunnerDomainTask.TaskTriggerType, triggeredByUserID *uuid.UUID) *uuid.UUID {
	switch triggerType {
	case taskRunnerDomainTask.TaskTriggerTypeSchedule,
		taskRunnerDomainTask.TaskTriggerTypeWebhook,
		taskRunnerDomainTask.TaskTriggerTypeAPI:
		return nil
	default:
		return triggeredByUserID
	}
}

func MergePayloadWithWorkflowNodes(prompt string, payloadNodes []dag.Node, workflowNodes []dag.Node) {
	if len(workflowNodes) == 0 {
		return
	}

	payloadNodesByID := make(map[string]dag.Node, len(payloadNodes))
	for _, payloadNode := range payloadNodes {
		payloadNodesByID[payloadNode.ID] = payloadNode
	}

	hasPrompt := strings.TrimSpace(prompt) != ""

	for i := range workflowNodes {
		if hasPrompt {
			workflowNodes[i].RuntimePrompt = prompt
		}

		payloadNode, exists := payloadNodesByID[workflowNodes[i].ID]
		if !exists {
			continue
		}

		if payloadNode.RuntimeInstruction != "" {
			workflowNodes[i].RuntimeInstruction = payloadNode.RuntimeInstruction
		}

		if payloadNode.Credentials != nil {
			workflowNodes[i].Credentials = payloadNode.Credentials
		}

		if payloadNode.Settings != nil {
			workflowNodes[i].Settings = payloadNode.Settings
		}
	}
}

func BuildNodeKeyWithBuilder(builder *strings.Builder, nodeType, nodeID string) string {
	builder.Reset()
	builder.Grow(len(nodeType) + len(nodeID) + 1)
	builder.WriteString(nodeType)
	builder.WriteString("_")
	builder.WriteString(nodeID)
	return builder.String()
}

func BuildNodeKey(nodeType, nodeID string) string {
	var builder strings.Builder
	return BuildNodeKeyWithBuilder(&builder, nodeType, nodeID)
}

func ApplyRuntimeConfig(rc *triggersContract.RuntimeConfig, workflowNodes []dag.Node) {
	if rc == nil {
		return
	}
	MergePayloadWithWorkflowNodes(rc.RuntimePrompt, rc.Nodes, workflowNodes)
}
