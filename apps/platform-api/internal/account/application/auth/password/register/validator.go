package register

import (
	"github.com/blocknextai/platform-api/internal/common/validation"
)

func (c *RegisterCommand) Validate() error {
	if err := validation.Email(c.Email); err != nil {
		return err
	}
	return validation.NewPassword(c.Password)
}
