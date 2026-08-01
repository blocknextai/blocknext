package taskrunner

import (
	"context"
	"log/slog"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/realtime"
	taskRunnerDomain "github.com/blocknextai/platform-api/internal/taskrunner/domain"
	taskRunnerDomainNode "github.com/blocknextai/platform-api/internal/taskrunner/domain/node"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
)

type EventPublisher struct {
	broadcaster realtime.Broadcaster
}

func NewEventPublisher(broadcaster realtime.Broadcaster) *EventPublisher {
	return &EventPublisher{
		broadcaster: broadcaster,
	}
}

func (p *EventPublisher) PublishTaskEvent(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	status taskRunnerDomain.Status,
	errorMsg string,
	duration int64,
) {
	event := taskRunnerDomainTask.NewTaskEvent(
		id,
		organizationID,
		executionContext,
		contextItemID,
		status,
		errorMsg,
		duration,
	)

	err := p.broadcaster.PublishTaskEvent(ctx, event)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to publish task event",
			"component", "event_publisher",
			"task_id", id,
			"error", err)
	}
}

func (p *EventPublisher) PublishNodeEvent(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	nodeID string,
	nodeType string,
	status taskRunnerDomain.Status,
	outputs []map[string]any,
	errorMsg string,
	duration int64,
) {
	event := taskRunnerDomainNode.NewNodeEvent(
		id,
		organizationID,
		executionContext,
		contextItemID,
		nodeID,
		nodeType,
		status,
		outputs,
		errorMsg,
		duration,
	)

	err := p.broadcaster.PublishNodeEvent(ctx, event)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to publish node event",
			"component", "event_publisher",
			"node_id", id,
			"error", err)
	}
}
