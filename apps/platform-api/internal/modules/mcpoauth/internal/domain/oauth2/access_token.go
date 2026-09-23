package oauth2

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID             uuid.UUID
	Subject        uuid.UUID
	OrganizationID uuid.UUID
	ClientID       string
	GrantID        uuid.UUID
	Scopes         Scopes
	Resource       string
	IssuedAt       time.Time
	ExpiresAt      time.Time
}
