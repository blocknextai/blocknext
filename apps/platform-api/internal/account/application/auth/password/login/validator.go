package login

import (
	"github.com/blocknextai/platform-api/internal/common/validation"
)

func (c *LoginCommand) Validate() error {
	if err := validation.Email(c.Email); err != nil {
		return err
	}
	return validation.Password(c.Password)
}
