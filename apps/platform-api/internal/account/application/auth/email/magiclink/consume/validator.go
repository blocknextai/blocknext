package consume

import (
	"strings"

	"github.com/blocknextai/go-packages/apperror"
)

var (
	errTokenIsRequired = apperror.Validation("token is required")
)

func (c *MagicLinkConsumeCommand) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return errTokenIsRequired
	}
	return nil
}
