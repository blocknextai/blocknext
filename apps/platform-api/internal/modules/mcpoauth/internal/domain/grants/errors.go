package grants

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrGrantNotFound         = apperror.NotFound("grant not found")
	ErrInvalidGrantID        = apperror.Validation("invalid grant id")
	ErrGrantRevoked          = apperror.Validation("grant revoked")
	ErrResourceMismatch      = apperror.Validation("resource does not match the resource the grant was issued for")
	ErrScopeExceedsGrant     = apperror.Validation("requested scope exceeds the granted scope")
	ErrInvalidClientID       = apperror.Validation("invalid client id")
	ErrInvalidUserID         = apperror.Validation("invalid user id")
	ErrInvalidOrganizationID = apperror.Validation("invalid organization id")
	ErrInvalidResource       = apperror.Validation("invalid resource")
	ErrInvalidScopes         = apperror.Validation("invalid scopes")
)
