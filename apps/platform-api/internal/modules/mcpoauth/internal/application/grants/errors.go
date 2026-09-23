package grants

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToCreateGrant = apperror.Internal("failed to create grant")
	ErrFailedToUpdateGrant = apperror.Internal("failed to update grant")
	ErrFailedToRevokeGrant = apperror.Internal("failed to revoke grant")
)
