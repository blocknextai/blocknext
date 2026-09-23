package generation

import (
	"context"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	generationDomainSessions "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/generation/sessions"
)

type GetAllSessionsQuery struct {
	OrganizationID uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}

type SessionResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetAllSessionsResponse struct {
	Items      []*SessionResponse
	TotalCount int64
}

func (s *Service) GetAllSessions(ctx context.Context, request *GetAllSessionsQuery) (*GetAllSessionsResponse, error) {
	sessions, totalCount, err := s.sessionRepository.GetAllByOrganizationID(
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

func MapGetAllSessionsQueryToGetAllSessionsResponse(sessions []*generationDomainSessions.GenerationSession, totalCount int64) *GetAllSessionsResponse {
	items := make([]*SessionResponse, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, &SessionResponse{
			ID:        session.ID,
			Title:     session.Title,
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
		})
	}

	return &GetAllSessionsResponse{
		Items:      items,
		TotalCount: totalCount,
	}
}
