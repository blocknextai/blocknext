package getallsessions

import (
	"context"

	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
)

type Handler struct {
	sessionRepository generationDomainSessions.SessionRepository
}

func New(
	sessionRepository generationDomainSessions.SessionRepository,
) *Handler {
	return &Handler{
		sessionRepository: sessionRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllSessionsQuery) (*GetAllSessionsResponse, error) {
	sessions, totalCount, err := h.sessionRepository.GetAllByOrganizationID(
		ctx,
		request.OrganizationID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	return MapGetAllSessionsQueryToGetAllSessionsResponse(sessions, totalCount), nil
}
