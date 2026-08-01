package updatesession

import (
	"time"

	"github.com/google/uuid"
)

type UpdateSessionResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}
