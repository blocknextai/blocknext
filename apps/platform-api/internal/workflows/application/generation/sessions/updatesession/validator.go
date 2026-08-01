package updatesession

import (
	"strings"

	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
)

func (c *UpdateSessionCommand) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return generationDomainSessions.ErrTitleIsRequired
	}

	return nil
}
