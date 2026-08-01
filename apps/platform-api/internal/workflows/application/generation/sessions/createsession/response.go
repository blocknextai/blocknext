package createsession

import (
	"time"

	"github.com/google/uuid"
)

type CreateSessionResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}
