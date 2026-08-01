package organizationusers

import (
	"github.com/google/uuid"
)

type OrganizationUserRoleChangedDomainEvent struct {
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	OldRole        string    `json:"oldRole"`
	NewRole        string    `json:"newRole"`
}

func (OrganizationUserRoleChangedDomainEvent) EventName() string {
	return "organizations.organization_user.role_changed"
}
