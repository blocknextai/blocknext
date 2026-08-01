package organizations

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToUpdateOrganization = apperror.Internal("failed to update organization")
	ErrFailedToDeleteOrganization = apperror.Internal("failed to delete organization")
	ErrInvalidTitle               = apperror.Validation("invalid title")
	ErrTitleTooLong               = apperror.Validation("title too long")
)
