package organizations

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrOrganizationNotFound = apperror.NotFound("organization not found")
	ErrTitleIsRequired      = apperror.Validation("title is required")
)
