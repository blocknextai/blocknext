package createsession

import (
	"strings"

	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
)

func (c *CreateSessionCommand) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return generationDomainSessions.ErrTitleIsRequired
	}

	return nil
}
