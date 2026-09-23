package history

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
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
	toolInvocationService executionsContract.ToolInvocationService
}

func NewRecorder(toolInvocationService executionsContract.ToolInvocationService) Recorder {
	return &recorder{
		toolInvocationService: toolInvocationService,
	}
}

func (r *recorder) Record(ctx context.Context, call ToolCall) {
	status := executionsContract.StatusSuccess
	var errorMessage *string
	if call.Err != nil {
		status = executionsContract.StatusFailed
		errorMessage = new(call.Err.Error())
	}

	if _, err := r.toolInvocationService.Record(
		ctx,
		call.OrganizationID,
		call.APIKeyID,
		executionsContract.SourceMCP,
		call.ToolID,
		status,
		call.Parameters,
		call.Credentials,
		call.Outputs,
		errorMessage,
		call.StartedAt,
		call.CompletedAt,
	); err != nil {
		slog.WarnContext(ctx, "failed to record mcp tool call",
			"component", "mcp_history_recorder",
			"organization_id", call.OrganizationID,
			"tool_id", call.ToolID,
			"error", err)
	}
}
