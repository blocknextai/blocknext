package add

import (
	"github.com/google/uuid"
)

type AddEmailCommand struct {
	UserID uuid.UUID
	Email  string
}
