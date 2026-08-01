package set

import (
	"github.com/google/uuid"
)

type SetPasswordCommand struct {
	UserID   uuid.UUID
	Password string
}
