package functioncalling

import (
	"context"
)

type FunctionCallingService interface {
	ExecuteWithContext(ctx context.Context, data []map[string]any, functionDeclarations []map[string]any) ([]map[string]any, error)
}

type FunctionCallingInput struct {
	Instruction        string `json:"instruction"`
	RuntimeInstruction string `json:"runtimeInstruction"`
	RuntimePrompt      string `json:"runtimePrompt"`
}
