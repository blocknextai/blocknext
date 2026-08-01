package authurl

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/google/uuid"
)

var ErrCredentialIDIsRequired = apperror.Validation("credential id is required")

func (c *AuthURLCommand) Validate() error {
	if c.CredentialID == uuid.Nil {
		return ErrCredentialIDIsRequired
	}
	return nil
}
