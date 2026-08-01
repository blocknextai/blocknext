package change

import (
	"github.com/google/uuid"
)

type ChangeEmailCommand struct {
	UserID          uuid.UUID
	NewEmail        string
	CurrentPassword string
}
