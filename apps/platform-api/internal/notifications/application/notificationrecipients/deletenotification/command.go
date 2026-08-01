package deletenotification

import (
	"github.com/google/uuid"
)

type DeleteNotificationCommand struct {
	UserID      uuid.UUID
	RecipientID uuid.UUID
}
