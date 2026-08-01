package change

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/common/validation"
)

var ErrSameAsCurrent = apperror.Validation("new password must differ from current password")

func (c *ChangePasswordCommand) Validate() error {
	if err := validation.Password(c.CurrentPassword); err != nil {
		return err
	}
	if err := validation.NewPassword(c.NewPassword); err != nil {
		return err
	}
	if c.CurrentPassword == c.NewPassword {
		return ErrSameAsCurrent
	}
	return nil
}
