package history

import (
	"context"
	"log/slog"
	"time"

	executionsApplicationToolInvocations "github.com/blocknextai/platform-api/internal/executions/application/toolinvocations"
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
	"github.com/blocknextai/platform-api/internal/realtime"
	"github.com/google/uuid"
)

type ToolCall struct {
	OrganizationID uuid.UUID
	APIKeyID       *uuid.UUID
	ToolID         string
	Parameters     map[string]any
	Credentials    map[string]any
	Outputs        []map[string]any
	Err            error
	StartedAt      time.Time
	CompletedAt    time.Time
}

type Recorder interface {
	Record(ctx context.Context, call ToolCall)
}

type recorder struct {
	toolInvocationService executionsApplicationToolInvocations.ToolInvocationService
	broadcaster           realtime.Broadcaster
}

func NewRecorder(
	toolInvocationService executionsApplicationToolInvocations.ToolInvocationService,
	broadcaster realtime.Broadcaster,
) Recorder {
	return &recorder{
		toolInvocationService: toolInvocationService,
		broadcaster:           broadcaster,
	}
}

func (r *recorder) Record(ctx context.Context, call ToolCall) {
	status := executionsDomainToolInvocations.StatusSuccess
	errorText := ""
	var errorMessage *string
	if call.Err != nil {
		status = executionsDomainToolInvocations.StatusFailed
		errorText = call.Err.Error()
		errorMessage = &errorText
	}

	id, err := r.toolInvocationService.Record(
		ctx,
		call.OrganizationID,
		call.APIKeyID,
		executionsDomainToolInvocations.SourceMCP,
		call.ToolID,
		status,
		call.Parameters,
		call.Credentials,
		call.Outputs,
		errorMessage,
		call.StartedAt,
		call.CompletedAt,
	)
	if err != nil {
		slog.WarnContext(ctx, "failed to record mcp tool call",
			"component", "mcp_history_recorder",
			"organization_id", call.OrganizationID,
			"tool_id", call.ToolID,
			"error", err)
		return
	}

	r.broadcast(ctx, id, call.OrganizationID, call.ToolID, status, errorText)
}

func (r *recorder) broadcast(
	ctx context.Context,
	id uuid.UUID,
	organizationID uuid.UUID,
	toolID string,
	status executionsDomainToolInvocations.Status,
	errorMessage string,
) {
	if r.broadcaster == nil {
		return
	}

	event := executionsDomainToolInvocations.NewToolInvocationEvent(
		id,
		organizationID,
		executionsDomainToolInvocations.SourceMCP,
		toolID,
		status,
		errorMessage,
	)

	if err := r.broadcaster.PublishToolInvocationEvent(ctx, event); err != nil {
		slog.WarnContext(ctx, "failed to publish mcp tool invocation event",
			"component", "mcp_history_recorder",
			"organization_id", organizationID,
			"tool_id", toolID,
			"error", err)
	}
}
