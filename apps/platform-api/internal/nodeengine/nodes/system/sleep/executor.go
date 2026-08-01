package sleep

import (
	"context"
	"time"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
)

type SleepExecutorInput struct {
	Duration float64 `schema:"duration"`
}

type SleepExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[SleepExecutorInput]
}

func NewSleepExecutor(
	nodeID string,
	validator *jsonschema.Validator[SleepExecutorInput],
) *SleepExecutor {
	return &SleepExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

func (e *SleepExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	input, err := e.validator.Parse(data[0])
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(time.Duration(input.Duration) * time.Millisecond):
		return []map[string]any{
			{
				"status": true,
			},
		}, nil
	}
}
