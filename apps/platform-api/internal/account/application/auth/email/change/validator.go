package change

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/common/validation"
)

var ErrSameEmail = apperror.Validation("new email is the same as current email")

func (c *ChangeEmailCommand) Validate() error {
	if err := validation.Email(c.NewEmail); err != nil {
		return err
	}
	return validation.Password(c.CurrentPassword)
}
