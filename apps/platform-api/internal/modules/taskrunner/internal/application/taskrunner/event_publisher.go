package taskrunner

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	taskRunnerDomain "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain"
	"github.com/blocknextai/platform-api/internal/realtime"
	"github.com/blocknextai/platform-api/internal/realtime/events"
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
	event := &events.TaskEvent{
		ID:               id,
		Type:             events.TaskEventType,
		OrganizationID:   organizationID,
		ExecutionContext: executionContext,
		ContextItemID:    contextItemID,
		Status:           string(status),
		Error:            errorMsg,
		Duration:         duration,
	}

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
	event := &events.NodeEvent{
		ID:               id,
		Type:             events.NodeEventType,
		OrganizationID:   organizationID,
		ExecutionContext: executionContext,
		ContextItemID:    contextItemID,
		NodeID:           nodeID,
		NodeType:         nodeType,
		Status:           string(status),
		Outputs:          outputs,
		Error:            errorMsg,
		Duration:         duration,
	}

	err := p.broadcaster.PublishNodeEvent(ctx, event)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to publish node event",
			"component", "event_publisher",
			"node_id", id,
			"error", err)
	}
}
