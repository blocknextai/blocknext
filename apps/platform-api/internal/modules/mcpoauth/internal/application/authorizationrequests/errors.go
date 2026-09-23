package authorizationrequests

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToCreateAuthorizationRequest  = apperror.Internal("failed to create authorization request")
	ErrFailedToApproveAuthorizationRequest = apperror.Internal("failed to approve authorization request")
	ErrFailedToDenyAuthorizationRequest    = apperror.Internal("failed to deny authorization request")
)
