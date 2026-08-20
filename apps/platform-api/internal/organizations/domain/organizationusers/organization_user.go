package organizationusers

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/rbac"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type OrganizationUser struct {
	database.BaseEntity

	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	Alias          string
}

func New(
	organizationID uuid.UUID,
	userID uuid.UUID,
	role string,
	alias string,
) (*OrganizationUser, error) {
	utcNow := time.Now().UTC()

	organizationUser := &OrganizationUser{
		ID:             bnuuid.NewV7(),
		CreatedAt:      utcNow,
		UpdatedAt:      utcNow,
		DeletedAt:      nil,
		OrganizationID: organizationID,
		UserID:         userID,
		Role:           role,
		Alias:          alias,
	}

	return organizationUser.validateThenReturn()
}

func (o *OrganizationUser) Update(
	role string,
	alias string,
) (*OrganizationUser, error) {
	o.UpdatedAt = time.Now().UTC()

	o.Role = role
	o.Alias = alias

	return o.validateThenReturn()
}

func (o *OrganizationUser) Delete() (*OrganizationUser, error) {
	utcNow := time.Now().UTC()

	o.UpdatedAt = utcNow
	o.DeletedAt = new(utcNow)
	return o.validateThenReturn()
}

func (o *OrganizationUser) validateThenReturn() (*OrganizationUser, error) {
	if o.OrganizationID == uuid.Nil {
		return nil, ErrOrganizationIDIsRequired
	}

	if o.UserID == uuid.Nil {
		return nil, ErrUserIDIsRequired
	}

	if strings.TrimSpace(o.Role) == "" {
		return nil, ErrRoleIsRequired
	}

	if !rbac.IsValidOrganizationRole(o.Role) {
		return nil, ErrInvalidRole
	}

	return o, nil
}
