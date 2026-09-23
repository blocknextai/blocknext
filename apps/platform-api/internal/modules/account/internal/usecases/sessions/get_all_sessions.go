package sessions

import (
	"context"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/sessions"
)

type GetAllSessionsQuery struct {
	UserID     uuid.UUID
	SessionID  uuid.UUID
	Search     resultPkg.SearchRequest
	Pagination resultPkg.PaginationRequest
}

type SessionResponse struct {
	SessionID    uuid.UUID                  `json:"sessionId"`
	AuthProvider accountDomain.AuthProvider `json:"authProvider"`
	UserAgent    string                     `json:"userAgent"`
	CreatedAt    time.Time                  `json:"createdAt"`
	UpdatedAt    time.Time                  `json:"updatedAt"`
	IsCurrent    bool                       `json:"isCurrent"`
}

type GetAllSessionsResponse struct {
	Items      []SessionResponse `json:"items"`
	TotalCount int64             `json:"totalCount"`
}

func (s *Service) GetAllSessions(ctx context.Context, request *GetAllSessionsQuery) (*GetAllSessionsResponse, error) {
	sessions, total, err := s.sessionRepository.GetActiveByUserID(
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

func MapSessionsToResponse(sessions []*sessions.Session, currentSessionID uuid.UUID) []SessionResponse {
	response := make([]SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		response = append(response, SessionResponse{
			SessionID:    s.ID,
			AuthProvider: s.AuthProvider,
			UserAgent:    s.UserAgent,
			CreatedAt:    s.CreatedAt,
			UpdatedAt:    s.UpdatedAt,
			IsCurrent:    s.ID == currentSessionID,
		})
	}
	return response
}
