package createorganization

import (
	"strings"

	organizationsApplicationOrganizations "github.com/blocknextai/platform-api/internal/organizations/application/organizations"
)

const (
	MaxTitleLength = 255
)

func (c *CreateOrganizationCommand) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return organizationsApplicationOrganizations.ErrInvalidTitle
	}

	if len(strings.TrimSpace(c.Title)) > MaxTitleLength {
		return organizationsApplicationOrganizations.ErrTitleTooLong
	}

	return nil
}
