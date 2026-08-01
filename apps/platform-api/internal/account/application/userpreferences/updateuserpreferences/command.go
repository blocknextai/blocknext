package updateuserpreferences

import (
	"github.com/google/uuid"
)

type ThemeCommand struct {
	Color *string
	Mode  *string
}

type UpdateUserPreferencesCommand struct {
	UserID   uuid.UUID
	Theme    *ThemeCommand
	Language *string
}
