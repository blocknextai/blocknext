package starter

import (
	"context"

	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
)

type StarterExecutor struct {
	executors.Executor
}

func NewStarterExecutor(nodeID string) *StarterExecutor {
	return &StarterExecutor{
		ID: nodeID,
	}
}

func (e *StarterExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return []map[string]any{
			{
				"status": true,
			},
		}, nil
	}
}
