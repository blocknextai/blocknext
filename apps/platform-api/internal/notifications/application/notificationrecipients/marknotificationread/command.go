package marknotificationread

import (
	"github.com/google/uuid"
)

type MarkNotificationReadCommand struct {
	UserID      uuid.UUID
	RecipientID uuid.UUID
}
