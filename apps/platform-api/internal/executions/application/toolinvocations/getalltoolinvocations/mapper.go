package getalltoolinvocations

import (
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
)

func MapToolInvocationsToGetAllToolInvocationsResponse(
	toolInvocations []*executionsDomainToolInvocations.ToolInvocation,
) []*ToolInvocationResponse {
	items := make([]*ToolInvocationResponse, 0, len(toolInvocations))
	for _, toolInvocation := range toolInvocations {
		items = append(items, &ToolInvocationResponse{
			ID:           toolInvocation.ID,
			Source:       toolInvocation.Source.String(),
			ToolID:       toolInvocation.ToolID,
			Status:       toolInvocation.Status.String(),
			ErrorMessage: toolInvocation.ErrorMessage,
			StartedAt:    new(toolInvocation.StartedAt),
			CompletedAt:  toolInvocation.CompletedAt,
		})
	}

	return items
}
