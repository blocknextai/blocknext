package features

import (
	"context"
)

type GetFeaturesQuery struct{}

type GetFeaturesResponse struct {
	FunctionCalling     bool `json:"functionCalling"`
	WorkflowsGeneration bool `json:"workflowsGeneration"`
}

func (s *Service) GetFeatures(ctx context.Context, _ *GetFeaturesQuery) (*GetFeaturesResponse, error) {
	return &GetFeaturesResponse{
		FunctionCalling:     s.functionCallingEnabled,
		WorkflowsGeneration: s.workflowsGenerationEnabled,
	}, nil
}
