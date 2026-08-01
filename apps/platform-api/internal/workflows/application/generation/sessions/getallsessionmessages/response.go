package getallsessionmessages

import (
	"time"

	"github.com/google/uuid"
)

type MessageResponse struct {
	ID        uuid.UUID      `json:"id"`
	SessionID uuid.UUID      `json:"sessionId"`
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}

type GetAllSessionMessagesResponse struct {
	Items      []*MessageResponse
	TotalCount int64
}
