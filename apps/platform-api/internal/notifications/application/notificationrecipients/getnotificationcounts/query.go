package getnotificationcounts

import (
	"github.com/google/uuid"
)

type GetNotificationCountsQuery struct {
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
}
