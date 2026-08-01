package createcredential

import (
	"strings"

	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
)

const (
	MaxKeyLength   = 255
	MaxTitleLength = 255
)

func (c *CreateCredentialCommand) Validate() error {
	if !c.SourceType.IsValid() {
		return credentialsDomainCredentials.ErrInvalidSourceType
	}

	if strings.TrimSpace(c.Key) == "" {
		return credentialsDomainCredentials.ErrInvalidKey
	}

	if len(strings.TrimSpace(c.Key)) > MaxKeyLength {
		return credentialsDomainCredentials.ErrKeyTooLong
	}

	if strings.TrimSpace(c.Title) == "" {
		return credentialsDomainCredentials.ErrInvalidTitle
	}

	if len(strings.TrimSpace(c.Title)) > MaxTitleLength {
		return credentialsDomainCredentials.ErrTitleTooLong
	}

	if c.SourceType == credentialsDomainCredentials.SourceTypeOwner {
		if len(c.Data) == 0 {
			return credentialsDomainCredentials.ErrInvalidData
		}
	}

	return nil
}
