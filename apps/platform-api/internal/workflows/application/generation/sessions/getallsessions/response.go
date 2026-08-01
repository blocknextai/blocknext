package getallsessions

import (
	"time"

	"github.com/google/uuid"
)

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
