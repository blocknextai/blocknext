package forgot

import (
	"github.com/blocknextai/platform-api/internal/common/validation"
)

func (c *ForgotCommand) Validate() error {
	return validation.Email(c.Email)
}
