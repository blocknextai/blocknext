package passwordpolicy

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrPasswordBreached = apperror.Validation("password appears in known data breaches; choose a different one")
)
