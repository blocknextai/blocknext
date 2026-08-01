package nodeexecutions

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type NodeExecution struct {
	database.BaseEntity

	TaskID                 uuid.UUID
	NodeType               string
	NodeID                 string
	Status                 string
	Inputs                 []map[string]any
	Outputs                []map[string]any
	FunctionCallingOutputs []map[string]any
	ErrorMessage           *string
	StartedAt              *time.Time
	CompletedAt            *time.Time
}

func New(
	taskID uuid.UUID,
	nodeType string,
	nodeID string,
	status string,
	inputs []map[string]any,
	errorMessage *string,
	startedAt *time.Time,
) (*NodeExecution, error) {
	utcNow := time.Now().UTC()

	nodeExecution := &NodeExecution{
		BaseEntity: database.BaseEntity{
			ID:        bnuuid.NewV7(),
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		TaskID:                 taskID,
		NodeType:               nodeType,
		NodeID:                 nodeID,
		Status:                 status,
		Inputs:                 inputs,
		Outputs:                nil,
		FunctionCallingOutputs: nil,
		ErrorMessage:           errorMessage,
		StartedAt:              startedAt,
	}

	return nodeExecution.validateThenReturn()
}

func (n *NodeExecution) Update(
	status string,
	inputs []map[string]any,
	outputs []map[string]any,
	functionCallingOutputs []map[string]any,
	errorMessage *string,
	startedAt *time.Time,
	completedAt *time.Time,
) (*NodeExecution, error) {
	n.UpdatedAt = time.Now().UTC()

	n.Status = status
	n.Inputs = inputs
	n.Outputs = outputs
	n.FunctionCallingOutputs = functionCallingOutputs
	n.ErrorMessage = errorMessage
	n.StartedAt = startedAt
	n.CompletedAt = completedAt

	return n.validateThenReturn()
}

func (w *NodeExecution) validateThenReturn() (*NodeExecution, error) {
	if w.TaskID == uuid.Nil {
		return nil, ErrTaskIDIsRequired
	}

	if strings.TrimSpace(w.NodeType) == "" {
		return nil, ErrNodeTypeIsRequired
	}

	if strings.TrimSpace(w.NodeID) == "" {
		return nil, ErrNodeIDIsRequired
	}

	if strings.TrimSpace(w.Status) == "" {
		return nil, ErrStatusIsRequired
	}

	return w, nil
}
