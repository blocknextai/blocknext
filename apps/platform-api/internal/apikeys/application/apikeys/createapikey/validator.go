package createapikey

import (
	"strings"

	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
)

func (c *CreateAPIKeyCommand) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return apiKeysDomainAPIKeys.ErrInvalidAPIKeyName
	}

	if len(c.Scopes) == 0 {
		return apiKeysDomainAPIKeys.ErrInvalidAPIKeyScopes
	}
	for _, scope := range c.Scopes {
		if !scope.IsValid() {
			return apiKeysDomainAPIKeys.ErrInvalidAPIKeyScopes
		}
	}

	return nil
}
