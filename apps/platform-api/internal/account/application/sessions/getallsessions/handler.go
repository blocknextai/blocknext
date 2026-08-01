package getallsessions

import (
	"context"

	accountDomainSessions "github.com/blocknextai/platform-api/internal/account/domain/sessions"
)

type Handler struct {
	sessionRepository accountDomainSessions.SessionRepository
}

func New(
	sessionRepository accountDomainSessions.SessionRepository,
) *Handler {
	return &Handler{
		sessionRepository: sessionRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllSessionsQuery) (*GetAllSessionsResponse, error) {
	sessions, total, err := h.sessionRepository.GetActiveByUserID(
		ctx,
		request.UserID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	return &GetAllSessionsResponse{
		Items:      MapSessionsToResponse(sessions, request.SessionID),
		TotalCount: total,
	}, nil
}
