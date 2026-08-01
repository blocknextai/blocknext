package createorganizationuser

import (
	"strings"

	organizationsApplicationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
)

func (c *CreateOrganizationUserCommand) Validate() error {
	if strings.TrimSpace(c.Identifier) == "" {
		return organizationsApplicationUsers.ErrIdentifierIsRequired
	}

	return nil
}
