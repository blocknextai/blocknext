package triggertask

import (
	"context"
	"strings"

	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	taskRunnerApplicationContextresolver "github.com/blocknextai/platform-api/internal/taskrunner/application/contextresolver"
	taskRunnerApplication "github.com/blocknextai/platform-api/internal/taskrunner/application/taskrunner"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
)

type TriggerTaskHandler struct {
	taskService     taskRunnerDomainTaskRunner.TaskService
	contextResolver taskRunnerApplicationContextresolver.ContextResolver
}

func NewTriggerTaskHandler(
	taskService taskRunnerDomainTaskRunner.TaskService,
	contextResolver taskRunnerApplicationContextresolver.ContextResolver,
) *TriggerTaskHandler {
	return &TriggerTaskHandler{
		taskService:     taskService,
		contextResolver: contextResolver,
	}
}

func (h *TriggerTaskHandler) Handle(ctx context.Context, command *TriggerTaskCommand) (*TriggerTaskResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	resolvedContext, err := h.contextResolver.ResolveContext(ctx, command.ExecutionContext, command.ContextItemID, command.OrganizationID)
	if err != nil {
		return nil, err
	}

	var runtimeConfig *triggersDomainTriggers.RuntimeConfig
	if command.TriggerType == taskRunnerDomainTask.TaskTriggerTypeSchedule ||
		command.TriggerType == taskRunnerDomainTask.TaskTriggerTypeWebhook {
		runtimeConfig = &triggersDomainTriggers.RuntimeConfig{
			Nodes: command.Nodes,
		}
		if command.TriggerType == taskRunnerDomainTask.TaskTriggerTypeSchedule {
			runtimeConfig.RuntimePrompt = command.RuntimePrompt
		}
	}

	taskRunnerApplication.MergePayloadWithWorkflowNodes(command.RuntimePrompt, command.Nodes, resolvedContext.Nodes)

	var triggerContext *nodeEngineDomainAdapters.TriggerContext
	if strings.TrimSpace(command.RuntimePrompt) != "" {
		sender := ""
		if command.TriggeredByUserID != nil {
			sender = command.TriggeredByUserID.String()
		}
		triggerContext = &nodeEngineDomainAdapters.TriggerContext{
			Source: command.TriggerType.String(),
			Sender: sender,
			Prompt: command.RuntimePrompt,
		}
	}

	task, err := h.taskService.ExecuteTask(
		ctx,
		command.TriggeredByUserID,
		command.OrganizationID,
		command.ExecutionContext,
		command.ContextItemID,
		command.TriggerType,
		triggerContext,
		command.CronPattern,
		runtimeConfig,
		resolvedContext.Nodes,
		resolvedContext.Edges,
	)
	if err != nil {
		return nil, err
	}

	return &TriggerTaskResponse{
		ID:           task.ID,
		WebhookToken: task.WebhookToken,
	}, nil
}
