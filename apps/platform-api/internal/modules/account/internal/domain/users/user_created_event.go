package users

import (
	"github.com/google/uuid"
)

type UserCreatedDomainEvent struct {
	UserID      uuid.UUID `json:"userId"`
	Identifier  string    `json:"identifier"`
	DisplayName string    `json:"displayName"`
}

func (UserCreatedDomainEvent) EventName() string {
	return "account.user.created"
}
