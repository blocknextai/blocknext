package nodeexecutions

import (
	"context"

	"github.com/google/uuid"
)

type NodeExecutionRepository interface {
	GetByIDAndTaskID(ctx context.Context, nodeExecutionID uuid.UUID, taskID uuid.UUID) (*NodeExecution, error)
	GetAllByTaskID(ctx context.Context, taskID uuid.UUID) ([]*NodeExecution, error)
	Create(ctx context.Context, nodeExecution *NodeExecution) error
	Update(ctx context.Context, nodeExecution *NodeExecution) error
}
