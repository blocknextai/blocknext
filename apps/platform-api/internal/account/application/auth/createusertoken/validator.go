package createusertoken

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidPayload = apperror.Validation("invalid payload")
)

func (c *CreateUserTokenCommand) Validate() error {
	if c.Payload == nil {
		return ErrInvalidPayload
	}

	return nil
}
