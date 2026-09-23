package events

import (
	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/json"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
)

const (
	SubscriberBuffer = 64

	TaskEventType           = "task"
	NodeEventType           = "node"
	ToolInvocationEventType = "tool_invocation"
)

type TaskEvent struct {
	ID               uuid.UUID                     `json:"id,omitempty"`
	Type             string                        `json:"type,omitempty"`
	OrganizationID   uuid.UUID                     `json:"organizationId,omitempty"`
	ExecutionContext commonDomain.ExecutionContext `json:"executionContext,omitempty"`
	ContextItemID    uuid.UUID                     `json:"contextItemId,omitempty"`
	Status           string                        `json:"status,omitempty"`
	Error            string                        `json:"error,omitempty"`
	Duration         int64                         `json:"duration,omitempty"`
}

type NodeEvent struct {
	ID               uuid.UUID                     `json:"id,omitempty"`
	Type             string                        `json:"type,omitempty"`
	OrganizationID   uuid.UUID                     `json:"organizationId,omitempty"`
	ExecutionContext commonDomain.ExecutionContext `json:"executionContext,omitempty"`
	ContextItemID    uuid.UUID                     `json:"contextItemId,omitempty"`
	NodeID           string                        `json:"nodeId,omitempty"`
	NodeType         string                        `json:"nodeType,omitempty"`
	Status           string                        `json:"status,omitempty"`
	Outputs          []map[string]any              `json:"outputs,omitempty"`
	Error            string                        `json:"error,omitempty"`
	Duration         int64                         `json:"duration,omitempty"`
}

type ToolInvocationEvent struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Type           string    `json:"type,omitempty"`
	OrganizationID uuid.UUID `json:"organizationId,omitempty"`
	Source         string    `json:"source,omitempty"`
	ToolID         string    `json:"toolId,omitempty"`
	Status         string    `json:"status,omitempty"`
	Error          string    `json:"error,omitempty"`
}

func MarshalTask(event *TaskEvent) (string, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return "", ErrFailedToMarshalTaskEvent
	}

	return string(payload), nil
}

func MarshalNode(event *NodeEvent) (string, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return "", ErrFailedToMarshalNodeEvent
	}

	return string(payload), nil
}

func MarshalToolInvocation(event *ToolInvocationEvent) (string, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return "", ErrFailedToMarshalToolInvocationEvent
	}

	return string(payload), nil
}
