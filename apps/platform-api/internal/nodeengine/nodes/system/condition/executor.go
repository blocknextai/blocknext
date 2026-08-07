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
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

func (e *ConditionExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// TODO: only the first item is evaluated, so a node fed 10 items decides
		// the whole flow from item 0. Routing items per branch needs the runner
		// to carry item identity and keep outputs per edge; see the branching
		// notes before changing this.
		input, err := e.validator.Parse(data[0])
		if err != nil {
			return nil, err
		}

		result := false
		switch input.Operator {
		case "eq":
			result = input.LeftValue == input.RightValue
		case "neq":
			result = input.LeftValue != input.RightValue
		case "gt":
			result = cast.ToFloat(input.LeftValue) > cast.ToFloat(input.RightValue)
		case "gte":
			result = cast.ToFloat(input.LeftValue) >= cast.ToFloat(input.RightValue)
		case "lt":
			result = cast.ToFloat(input.LeftValue) < cast.ToFloat(input.RightValue)
		case "lte":
			result = cast.ToFloat(input.LeftValue) <= cast.ToFloat(input.RightValue)
		case "contains":
			result = strings.Contains(input.LeftValue, input.RightValue)
		case "not_contains":
			result = !strings.Contains(input.LeftValue, input.RightValue)
		case "is_empty":
			result = input.LeftValue == ""
		case "is_not_empty":
			result = input.LeftValue != ""
		default:
			return nil, ErrInvalidOperator
		}

		return []map[string]any{
			{
				"status": result,
			},
		}, nil
	}
}
