package deletenotification

import (
	"github.com/google/uuid"
)

type DeleteNotificationResponse struct {
	ID uuid.UUID `json:"id"`
}
