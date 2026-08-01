package task

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/taskrunner/domain"
	"github.com/google/uuid"
)

type TaskEvent struct {
	ID               uuid.UUID                     `json:"id,omitempty"`
	Type             string                        `json:"type,omitempty"`
	OrganizationID   uuid.UUID                     `json:"organizationId,omitempty"`
	ExecutionContext commonDomain.ExecutionContext `json:"executionContext,omitempty"`
	ContextItemID    uuid.UUID                     `json:"contextItemId,omitempty"`
	Status           domain.Status                 `json:"status,omitempty"`
	Error            string                        `json:"error,omitempty"`
	Duration         int64                         `json:"duration,omitempty"`
}

func NewTaskEvent(
	id uuid.UUID,
	organizationID uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	status domain.Status,
	errorMsg string,
	duration int64,
) *TaskEvent {
	return &TaskEvent{
		ID:               id,
		Type:             "task",
		OrganizationID:   organizationID,
		ExecutionContext: executionContext,
		ContextItemID:    contextItemID,
		Status:           status,
		Error:            errorMsg,
		Duration:         duration,
	}
}
