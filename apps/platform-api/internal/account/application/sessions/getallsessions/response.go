package getallsessions

import (
	"time"

	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/google/uuid"
)

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
