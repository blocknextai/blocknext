package users

import (
	"github.com/google/uuid"
)

type PasswordChangedDomainEvent struct {
	UserID uuid.UUID `json:"userId"`
}

func (PasswordChangedDomainEvent) EventName() string {
	return "account.user.password_changed"
}
