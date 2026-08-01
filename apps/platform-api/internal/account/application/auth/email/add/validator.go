package add

import (
	"github.com/blocknextai/platform-api/internal/common/validation"
)

func (c *AddEmailCommand) Validate() error {
	return validation.Email(c.Email)
}
