package markallnotificationsseen

import (
	"github.com/google/uuid"
)

type MarkAllNotificationsSeenCommand struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
}
