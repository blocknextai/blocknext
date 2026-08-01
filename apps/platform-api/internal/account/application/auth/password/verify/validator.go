package verify

import (
	"strings"

	"github.com/blocknextai/go-packages/apperror"
)

var (
	errTokenIsRequired = apperror.Validation("token is required")
)

func (c *VerifyCommand) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return errTokenIsRequired
	}
	return nil
}
