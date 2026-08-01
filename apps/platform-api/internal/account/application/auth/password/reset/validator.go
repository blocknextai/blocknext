package reset

import (
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/common/validation"
)

var (
	errTokenIsRequired = apperror.Validation("token is required")
)

func (c *ResetCommand) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return errTokenIsRequired
	}
	return validation.NewPassword(c.NewPassword)
}
