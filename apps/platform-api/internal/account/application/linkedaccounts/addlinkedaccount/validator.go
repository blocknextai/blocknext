package addlinkedaccount

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrPayloadRequired = apperror.Validation("payload is required")
)

func (c *AddLinkedAccountCommand) Validate() error {
	if c.Payload == nil {
		return ErrPayloadRequired
	}

	return nil
}
