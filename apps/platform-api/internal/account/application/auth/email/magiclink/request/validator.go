package request

import (
	"github.com/blocknextai/platform-api/internal/common/validation"
)

func (c *MagicLinkRequestCommand) Validate() error {
	return validation.Email(c.Email)
}
