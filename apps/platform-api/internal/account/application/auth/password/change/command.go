package change

import (
	"github.com/google/uuid"
)

type ChangePasswordCommand struct {
	UserID          uuid.UUID
	CurrentPassword string
	NewPassword     string
}
