package condition

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cast"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
)

var (
	ErrInvalidOperator = apperror.Internal("invalid operator")
)

type ConditionExecutorInput struct {
	LeftValue  string `schema:"leftValue"`
	Operator   string `schema:"operator"`
	RightValue string `schema:"rightValue"`
}

type ConditionExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[ConditionExecutorInput]
}

func NewConditionExecutor(
	nodeID string,
	validator *jsonschema.Validator[ConditionExecutorInput],
) *ConditionExecutor {
	return &ConditionExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

func (e *ConditionExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	outputs, _, err := e.ExecuteBranches(ctx, credentials, data)
	return outputs, err
}

func (e *ConditionExecutor) ExecuteBranches(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, map[string][]int, error) {
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
		outputs := make([]map[string]any, 0, len(data))
		branches := map[string][]int{
			BranchTrue:  make([]int, 0, len(data)),
			BranchFalse: make([]int, 0, len(data)),
		}

		for index, item := range data {
			result, err := e.evaluate(item)
			if err != nil {
				return nil, nil, err
			}

			outputs = append(outputs, map[string]any{"status": result})
			if result {
				branches[BranchTrue] = append(branches[BranchTrue], index)
				continue
			}
			branches[BranchFalse] = append(branches[BranchFalse], index)
		}

		return outputs, branches, nil
	}
}

func (e *ConditionExecutor) evaluate(item map[string]any) (bool, error) {
	input, err := e.validator.Parse(item)
	if err != nil {
		return false, err
	}

	left := strings.TrimSpace(input.LeftValue)
	right := strings.TrimSpace(input.RightValue)

	switch input.Operator {
	case "eq":
		return left == right, nil
	case "neq":
		return left != right, nil
	case "gt":
		return cast.ToFloat(left) > cast.ToFloat(right), nil
	case "gte":
		return cast.ToFloat(left) >= cast.ToFloat(right), nil
	case "lt":
		return cast.ToFloat(left) < cast.ToFloat(right), nil
	case "lte":
		return cast.ToFloat(left) <= cast.ToFloat(right), nil
	case "contains":
		return strings.Contains(left, right), nil
	case "not_contains":
		return !strings.Contains(left, right), nil
	case "is_empty":
		return left == "", nil
	case "is_not_empty":
		return left != "", nil
	default:
		return false, ErrInvalidOperator
	}
}
