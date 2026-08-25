package getalltoolinvocations

import (
	"context"

	"github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
)

type Handler struct {
	toolInvocationRepository toolinvocations.ToolInvocationRepository
}

func New(toolInvocationRepository toolinvocations.ToolInvocationRepository) *Handler {
	return &Handler{
		toolInvocationRepository: toolInvocationRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllToolInvocationsQuery) (*GetAllToolInvocationsResponse, error) {
	toolInvocations, totalCount, err := h.toolInvocationRepository.GetAllByOrganizationID(
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
