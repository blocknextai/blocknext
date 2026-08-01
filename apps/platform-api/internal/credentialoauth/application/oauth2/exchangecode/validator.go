package exchangecode

import (
	"strings"

	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrCodeIsRequired  = apperror.Validation("code is required")
	ErrStateIsRequired = apperror.Validation("state is required")
)

func (c *ExchangeCodeQuery) Validate() error {
	if strings.TrimSpace(c.Code) == "" {
		return ErrCodeIsRequired
	}

	if strings.TrimSpace(c.State) == "" {
		return ErrStateIsRequired
	}

	return nil
}
