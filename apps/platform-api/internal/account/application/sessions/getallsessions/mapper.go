package getallsessions

import (
	"github.com/blocknextai/platform-api/internal/account/domain/sessions"
	"github.com/google/uuid"
)

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
