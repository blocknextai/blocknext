package webhooks

import (
	"context"

	"github.com/google/uuid"

	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	taskRunnerApplicationContextResolver "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/contextresolver"
	taskRunnerApplication "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/taskrunner"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/taskrunner"
	triggersContract "github.com/blocknextai/platform-api/internal/modules/triggers/contract"
)

type ProcessResult struct {
	TaskID       *uuid.UUID                               `json:"taskId,omitempty"`
	Verification *nodeengineContract.VerificationResponse `json:"-"`
}

type WebhookProcessor interface {
	ProcessWebhook(ctx context.Context, req *triggersContract.Request) (*ProcessResult, error)
}

type webhookProcessor struct {
	webhookResolver triggersContract.WebhookResolver
	taskService     taskRunnerDomainTaskRunner.TaskService
	contextResolver taskRunnerApplicationContextResolver.ContextResolver
}

func NewWebhookProcessor(
	webhookResolver triggersContract.WebhookResolver,
	taskService taskRunnerDomainTaskRunner.TaskService,
	contextResolver taskRunnerApplicationContextResolver.ContextResolver,
) WebhookProcessor {
	return &webhookProcessor{
		webhookResolver: webhookResolver,
		taskService:     taskService,
		contextResolver: contextResolver,
	}
}

func (h *webhookProcessor) ProcessWebhook(ctx context.Context, req *triggersContract.Request) (*ProcessResult, error) {
	resolved, err := h.webhookResolver.Resolve(ctx, req)
	if err != nil {
		return nil, err
	}

	if resolved.Verification != nil {
		return &ProcessResult{Verification: resolved.Verification}, nil
	}

	resolvedContext, err := h.contextResolver.ResolveContext(ctx, resolved.ExecutionContext, resolved.ContextItemID, resolved.OrganizationID)
	if err != nil {
		return nil, err
	}

	taskRunnerApplication.ApplyRuntimeConfig(resolved.RuntimeConfig, resolvedContext.Nodes)
	taskRunnerApplication.MergePayloadWithWorkflowNodes(resolved.TriggerContext.Prompt, nil, resolvedContext.Nodes)

	task, err := h.taskService.ExecuteTask(
		ctx,
		resolved.TriggeredByUserID,
		resolved.OrganizationID,
		resolved.ExecutionContext,
		resolved.ContextItemID,
		taskRunnerDomainTask.TaskTriggerTypeWebhook,
		resolved.TriggerContext,
		nil,
		nil,
		resolvedContext.Nodes,
		resolvedContext.Edges,
	)
	if err != nil {
		return nil, err
	}

	return &ProcessResult{
		TaskID: &task.ID,
	}, nil
}
