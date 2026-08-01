package connections

import (
	"sync"

	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type Connection struct {
	ID             string
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Send           chan []byte
	closeOnce      sync.Once
}

func New(
	userID uuid.UUID,
	organizationID uuid.UUID,
) *Connection {
	return &Connection{
		ID:             bnuuid.NewV7().String(),
		UserID:         userID,
		OrganizationID: organizationID,
		Send:           make(chan []byte, 256),
	}
}

func (c *Connection) Close() {
	c.closeOnce.Do(func() {
		close(c.Send)
	})
}
