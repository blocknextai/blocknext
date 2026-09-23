package toolinvocations

import (
	"context"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/toolinvocations"
)

type GetAllToolInvocationsQuery struct {
	OrganizationID uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}

type ToolInvocationResponse struct {
	ID           uuid.UUID  `json:"id,omitempty"`
	Source       string     `json:"source,omitempty"`
	ToolID       string     `json:"toolId,omitempty"`
	Status       string     `json:"status,omitempty"`
	ErrorMessage *string    `json:"errorMessage,omitempty"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	CompletedAt  time.Time  `json:"completedAt"`
}

type GetAllToolInvocationsResponse struct {
	Items      []*ToolInvocationResponse `json:"items"`
	TotalCount int64                     `json:"totalCount"`
}

func (s *Service) GetAllToolInvocations(ctx context.Context, request *GetAllToolInvocationsQuery) (*GetAllToolInvocationsResponse, error) {
	toolInvocations, totalCount, err := s.toolInvocationRepository.GetAllByOrganizationID(
		ctx,
		request.OrganizationID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	return &GetAllToolInvocationsResponse{
		Items:      MapToolInvocationsToGetAllToolInvocationsResponse(toolInvocations),
		TotalCount: totalCount,
	}, nil
}

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
