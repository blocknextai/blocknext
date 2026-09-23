package authorizationrequests

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrAuthorizationRequestNotFound     = apperror.NotFound("authorization request not found")
	ErrAuthorizationRequestResolved     = apperror.Conflict("authorization request already resolved")
	ErrAuthorizationRequestExpired      = apperror.Validation("authorization request expired")
	ErrInvalidAuthorizationRequestID    = apperror.Validation("invalid authorization request id")
	ErrInvalidClientID                  = apperror.Validation("invalid client id")
	ErrInvalidRedirectURI               = apperror.Validation("invalid redirect uri")
	ErrUnsupportedResponseType          = apperror.Validation("only the authorization code response type is supported")
	ErrInvalidCodeChallenge             = apperror.Validation("code challenge is required")
	ErrInvalidCodeChallengeMethod       = apperror.Validation("code challenge method must be S256")
	ErrInvalidResource                  = apperror.Validation("unknown resource indicator")
	ErrInvalidScopes                    = apperror.Validation("invalid scopes")
	ErrInvalidStatus                    = apperror.Validation("invalid authorization request status")
	ErrInvalidUserID                    = apperror.Validation("invalid user id")
	ErrInvalidOrganizationID            = apperror.Validation("invalid organization id")
	ErrOrganizationMembershipRequired   = apperror.Forbidden("user is not a member of the organization")
	ErrAuthorizationRequestAccessDenied = apperror.Forbidden("the user denied the authorization request")
)
