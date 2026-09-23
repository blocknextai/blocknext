package tasks

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/dag"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	taskRunnerApplication "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/taskrunner"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/task"
	triggersContract "github.com/blocknextai/platform-api/internal/modules/triggers/contract"
)

type TriggerTaskCommand struct {
	TriggeredByUserID *uuid.UUID
	OrganizationID    uuid.UUID
	ExecutionContext  commonDomain.ExecutionContext
	ContextItemID     uuid.UUID
	TriggerType       taskRunnerDomainTask.TaskTriggerType
	CronPattern       *string

	RuntimePrompt string
	Nodes         []dag.Node
}

var (
	ErrInvalidExecutionContext = apperror.Validation("invalid execution context")
	ErrInvalidContextItemID    = apperror.Validation("invalid context item id")
	ErrInvalidTriggerType      = apperror.Validation("invalid trigger type")
	ErrInvalidNodes            = apperror.Validation("invalid nodes")
	ErrInvalidCronPattern      = apperror.Validation("invalid cron pattern")
)

func (c *TriggerTaskCommand) Validate() error {
	if !c.ExecutionContext.IsValid() {
		return ErrInvalidExecutionContext
	}

	if c.ContextItemID == uuid.Nil {
		return ErrInvalidContextItemID
	}

	if !c.TriggerType.IsValid() {
		return ErrInvalidTriggerType
	}

	if c.TriggerType == taskRunnerDomainTask.TaskTriggerTypeSchedule {
		if c.CronPattern == nil {
			return ErrInvalidCronPattern
		}
	}

	if len(c.Nodes) == 0 {
		return ErrInvalidNodes
	}

	return nil
}

type TriggerTaskResponse struct {
	ID           uuid.UUID `json:"id"`
	WebhookToken *string   `json:"webhookToken,omitempty"`
}

func (s *Service) TriggerTask(ctx context.Context, command *TriggerTaskCommand) (*TriggerTaskResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	resolvedContext, err := s.contextResolver.ResolveContext(ctx, command.ExecutionContext, command.ContextItemID, command.OrganizationID)
	if err != nil {
		return nil, err
	}

	var runtimeConfig *triggersContract.RuntimeConfig
	if command.TriggerType == taskRunnerDomainTask.TaskTriggerTypeSchedule ||
		command.TriggerType == taskRunnerDomainTask.TaskTriggerTypeWebhook {
		runtimeConfig = &triggersContract.RuntimeConfig{
			Nodes: command.Nodes,
		}
		if command.TriggerType == taskRunnerDomainTask.TaskTriggerTypeSchedule {
			runtimeConfig.RuntimePrompt = command.RuntimePrompt
		}
	}

	taskRunnerApplication.MergePayloadWithWorkflowNodes(command.RuntimePrompt, command.Nodes, resolvedContext.Nodes)

	var triggerContext *nodeengineContract.TriggerContext
	if strings.TrimSpace(command.RuntimePrompt) != "" {
		sender := ""
		if command.TriggeredByUserID != nil {
			sender = command.TriggeredByUserID.String()
		}
		triggerContext = &nodeengineContract.TriggerContext{
			Source: command.TriggerType.String(),
			Sender: sender,
			Prompt: command.RuntimePrompt,
		}
	}

	task, err := s.taskService.ExecuteTask(
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
