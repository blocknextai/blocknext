package getfeatures

import (
	"context"
)

type Handler struct {
	functionCallingEnabled     bool
	workflowsGenerationEnabled bool
}

func New(functionCallingEnabled bool, workflowsGenerationEnabled bool) *Handler {
	return &Handler{
		functionCallingEnabled:     functionCallingEnabled,
		workflowsGenerationEnabled: workflowsGenerationEnabled,
	}
}

func (h *Handler) Handle(_ context.Context, _ *GetFeaturesQuery) (*GetFeaturesResponse, error) {
	return &GetFeaturesResponse{
		FunctionCalling:     h.functionCallingEnabled,
		WorkflowsGeneration: h.workflowsGenerationEnabled,
	}, nil
}
