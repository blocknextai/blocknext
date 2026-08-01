package markallnotificationsread

import (
	"github.com/google/uuid"
)

type MarkAllNotificationsReadCommand struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
}
