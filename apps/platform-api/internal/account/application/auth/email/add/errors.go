package add

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrEmailAlreadyLinked = apperror.Conflict("email already linked to this account")
)
