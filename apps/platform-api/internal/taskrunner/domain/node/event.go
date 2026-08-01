package node

import (
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/taskrunner/domain"
	"github.com/google/uuid"
)

type NodeEvent struct {
	ID               uuid.UUID                     `json:"id,omitempty"`
	Type             string                        `json:"type,omitempty"`
	OrganizationID   uuid.UUID                     `json:"organizationId,omitempty"`
	ExecutionContext commonDomain.ExecutionContext `json:"executionContext,omitempty"`
	ContextItemID    uuid.UUID                     `json:"contextItemId,omitempty"`
	NodeID           string                        `json:"nodeId,omitempty"`
	NodeType         string                        `json:"nodeType,omitempty"`
	Status           domain.Status                 `json:"status,omitempty"`
	Outputs          []map[string]any              `json:"outputs,omitempty"`
	Error            string                        `json:"error,omitempty"`
	Duration         int64                         `json:"duration,omitempty"`
}

func NewNodeEvent(
	id uuid.UUID,
	organizationID uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	nodeID string,
	nodeType string,
	status domain.Status,
	outputs []map[string]any,
	errorMsg string,
	duration int64,
) *NodeEvent {
	return &NodeEvent{
		ID:               id,
		Type:             "node",
		OrganizationID:   organizationID,
		ExecutionContext: executionContext,
		ContextItemID:    contextItemID,
		NodeID:           nodeID,
		NodeType:         nodeType,
		Status:           status,
		Outputs:          outputs,
		Error:            errorMsg,
		Duration:         duration,
	}
}
