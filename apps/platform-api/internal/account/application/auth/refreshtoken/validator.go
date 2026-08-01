package refreshtoken

import (
	"strings"

	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrRefreshTokenRequired = apperror.Validation("refresh token is required")
)

func (c *RefreshTokenCommand) Validate() error {
	if strings.TrimSpace(c.RefreshToken) == "" {
		return ErrRefreshTokenRequired
	}

	return nil
}
