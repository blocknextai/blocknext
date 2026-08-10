package executors

import (
	"context"
)

type Executor struct {
	ID       string `json:"id,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

type BranchingExecutor interface {
	ExecuteBranches(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, map[string][]int, error)
}

type ExecutorManager interface {
	GetID() string
	GetDisabled() bool
	ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error)
}

func (e *Executor) GetID() string {
	return e.ID
}

func (e *Executor) GetDisabled() bool {
	return e.Disabled
}
