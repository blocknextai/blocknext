package taskexecutions

import (
	"time"

	"github.com/blocknextai/go-packages/dag"
	"github.com/blocknextai/go-packages/database"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type TaskExecution struct {
	database.BaseEntity

	OrganizationID    uuid.UUID
	TriggeredByUserID *uuid.UUID
	FlowTriggerID     *uuid.UUID
	ExecutionContext  commonDomain.ExecutionContext
	ContextItemID     uuid.UUID
	Status            string
	ExecutionType     ExecutionType
	ErrorMessage      *string
	Nodes             []dag.Node
	Edges             []dag.Edge
	StartedAt         *time.Time
	CompletedAt       *time.Time
}

func New(
	id uuid.UUID,
	organizationID uuid.UUID,
	triggeredByUserID *uuid.UUID,
	flowTriggerID *uuid.UUID,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	status string,
	executionType ExecutionType,
	errorMessage *string,
	nodes []dag.Node,
	edges []dag.Edge,
	startedAt *time.Time,
) (*TaskExecution, error) {
	utcNow := time.Now().UTC()

	taskExecution := &TaskExecution{
		BaseEntity: database.BaseEntity{
			ID:        id,
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		OrganizationID:    organizationID,
		TriggeredByUserID: triggeredByUserID,
		ExecutionContext:  executionContext,
		ContextItemID:     contextItemID,
		FlowTriggerID:     flowTriggerID,
		Status:            status,
		ExecutionType:     executionType,
		ErrorMessage:      errorMessage,
		Nodes:             nodes,
		Edges:             edges,
		StartedAt:         startedAt,
	}

	return taskExecution.validateThenReturn()
}

func (t *TaskExecution) Update(
	status string,
	errorMessage *string,
	completedAt *time.Time,
) (*TaskExecution, error) {
	t.UpdatedAt = time.Now().UTC()

	t.Status = status
	t.ErrorMessage = errorMessage
	t.CompletedAt = completedAt

	return t.validateThenReturn()
}

func (t *TaskExecution) Delete() (*TaskExecution, error) {
	utcNow := time.Now().UTC()

	t.UpdatedAt = utcNow
	t.DeletedAt = new(utcNow)
	return t.validateThenReturn()
}

func (t *TaskExecution) validateThenReturn() (*TaskExecution, error) {
	return t, nil
}
