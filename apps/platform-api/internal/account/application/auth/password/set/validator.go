package set

import (
	"github.com/blocknextai/platform-api/internal/common/validation"
)

func (c *SetPasswordCommand) Validate() error {
	return validation.NewPassword(c.Password)
}
