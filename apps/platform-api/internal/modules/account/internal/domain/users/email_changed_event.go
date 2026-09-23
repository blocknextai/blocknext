package users

import (
	"github.com/google/uuid"
)

type EmailChangedDomainEvent struct {
	UserID   uuid.UUID `json:"userId"`
	OldEmail string    `json:"oldEmail"`
	NewEmail string    `json:"newEmail"`
}

func (EmailChangedDomainEvent) EventName() string {
	return "account.user.email_changed"
}
