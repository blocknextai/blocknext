package webhooks

import (
	"context"

	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	taskRunnerApplicationContextResolver "github.com/blocknextai/platform-api/internal/taskrunner/application/contextresolver"
	taskRunnerApplication "github.com/blocknextai/platform-api/internal/taskrunner/application/taskrunner"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	triggersApplicationWebhooks "github.com/blocknextai/platform-api/internal/triggers/application/webhooks"
	"github.com/google/uuid"
)

type ProcessResult struct {
	TaskID       *uuid.UUID                                     `json:"taskId,omitempty"`
	Verification *nodeEngineDomainAdapters.VerificationResponse `json:"-"`
}

type WebhookProcessor interface {
	ProcessWebhook(ctx context.Context, req *triggersApplicationWebhooks.Request) (*ProcessResult, error)
}

type webhookProcessor struct {
	webhookResolver triggersApplicationWebhooks.WebhookResolver
	taskService     taskRunnerDomainTaskRunner.TaskService
	contextResolver taskRunnerApplicationContextResolver.ContextResolver
}

func NewWebhookProcessor(
	webhookResolver triggersApplicationWebhooks.WebhookResolver,
	taskService taskRunnerDomainTaskRunner.TaskService,
	contextResolver taskRunnerApplicationContextResolver.ContextResolver,
) WebhookProcessor {
	return &webhookProcessor{
		webhookResolver: webhookResolver,
		taskService:     taskService,
		contextResolver: contextResolver,
	}
}

func (h *webhookProcessor) ProcessWebhook(ctx context.Context, req *triggersApplicationWebhooks.Request) (*ProcessResult, error) {
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
