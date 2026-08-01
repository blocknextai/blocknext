package updatetrigger

import (
	"github.com/google/uuid"
)

type UpdateTriggerResponse struct {
	ID uuid.UUID `json:"id"`
}
