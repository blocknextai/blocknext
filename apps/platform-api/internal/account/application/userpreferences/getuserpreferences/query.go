package getuserpreferences

import (
	"github.com/google/uuid"
)

type GetUserPreferencesQuery struct {
	UserID uuid.UUID
}
