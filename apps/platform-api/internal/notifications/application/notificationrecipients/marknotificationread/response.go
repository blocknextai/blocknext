package marknotificationread

import (
	"github.com/google/uuid"
)

type MarkNotificationReadResponse struct {
	ID uuid.UUID `json:"id"`
}
