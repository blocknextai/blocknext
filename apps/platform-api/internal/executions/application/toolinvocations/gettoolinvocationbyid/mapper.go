package gettoolinvocationbyid

import (
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
)

func MapToolInvocationToResponse(
	toolInvocation *executionsDomainToolInvocations.ToolInvocation,
	apiKeyName *string,
) *GetToolInvocationByIDResponse {
	return &GetToolInvocationByIDResponse{
		ID:           toolInvocation.ID,
		Source:       toolInvocation.Source.String(),
		APIKeyID:     toolInvocation.APIKeyID,
		APIKeyName:   apiKeyName,
		ToolID:       toolInvocation.ToolID,
		Status:       toolInvocation.Status.String(),
		Parameters:   toolInvocation.Parameters,
		Outputs:      toolInvocation.Outputs,
		ErrorMessage: toolInvocation.ErrorMessage,
		StartedAt:    new(toolInvocation.StartedAt),
		CompletedAt:  toolInvocation.CompletedAt,
	}
}
