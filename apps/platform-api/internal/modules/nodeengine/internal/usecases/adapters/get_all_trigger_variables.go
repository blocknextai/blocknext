package adapters

import (
	"context"
)

type GetAllTriggerVariablesQuery struct{}

type GetAllTriggerVariablesResponse = []string

func (s *Service) GetAllTriggerVariables(ctx context.Context, _ *GetAllTriggerVariablesQuery) (*GetAllTriggerVariablesResponse, error) {
	triggerVariables := s.adapterService.GetTriggerVariables()
	return MapGetAllTriggerVariablesQueryToGetAllTriggerVariablesResponse(triggerVariables), nil
}

func MapGetAllTriggerVariablesQueryToGetAllTriggerVariablesResponse(
	variables []string,
) *GetAllTriggerVariablesResponse {
	response := make(GetAllTriggerVariablesResponse, 0, len(variables))
	response = append(response, variables...)
	return &response
}
