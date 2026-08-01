package deletelinkedaccount

import (
	"github.com/google/uuid"
)

type DeleteLinkedAccountCommand struct {
	UserID          uuid.UUID
	LinkedAccountID uuid.UUID
}
