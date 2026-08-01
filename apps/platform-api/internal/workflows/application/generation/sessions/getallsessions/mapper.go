package getallsessions

import (
	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
)

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
