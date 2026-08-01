package logout

import (
	"github.com/google/uuid"
)

type LogoutCommand struct {
	SessionID uuid.UUID
}
