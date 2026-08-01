package resendverification

import (
	"github.com/blocknextai/platform-api/internal/common/validation"
)

func (c *ResendVerificationCommand) Validate() error {
	return validation.Email(c.Email)
}
