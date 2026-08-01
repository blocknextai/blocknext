package canceltask

import (
	"github.com/google/uuid"
)

type CancelTaskResponse struct {
	ID uuid.UUID `json:"id"`
}
