package organizationusers

import (
	"github.com/google/uuid"
)

type OrganizationUserCreatedDomainEvent struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
}

func (OrganizationUserCreatedDomainEvent) EventName() string {
	return "organizations.organization_user.created"
}
