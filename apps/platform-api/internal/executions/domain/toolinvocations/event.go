package toolinvocations

import (
	"github.com/google/uuid"
)

const EventType = "tool_invocation"

type ToolInvocationEvent struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Type           string    `json:"type,omitempty"`
	OrganizationID uuid.UUID `json:"organizationId,omitempty"`
	Source         Source    `json:"source,omitempty"`
	ToolID         string    `json:"toolId,omitempty"`
	Status         Status    `json:"status,omitempty"`
	Error          string    `json:"error,omitempty"`
}

func NewToolInvocationEvent(
	id uuid.UUID,
	organizationID uuid.UUID,
	source Source,
	toolID string,
	status Status,
	errorMsg string,
) *ToolInvocationEvent {
	return &ToolInvocationEvent{
		ID:             id,
		Type:           EventType,
		OrganizationID: organizationID,
		Source:         source,
		ToolID:         toolID,
		Status:         status,
		Error:          errorMsg,
	}
}
