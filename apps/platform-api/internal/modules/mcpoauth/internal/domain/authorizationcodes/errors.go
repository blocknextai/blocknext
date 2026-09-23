package authorizationcodes

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrAuthorizationCodeNotFound         = apperror.NotFound("authorization code not found")
	ErrAuthorizationCodeUsed             = apperror.Validation("authorization code already used")
	ErrAuthorizationCodeExpired          = apperror.Validation("authorization code expired")
	ErrAuthorizationCodeClientMismatch   = apperror.Validation("authorization code was issued to another client")
	ErrAuthorizationCodeRedirectMismatch = apperror.Validation("redirect uri does not match the authorization request")
	ErrInvalidCodeVerifier               = apperror.Validation("invalid code verifier")
	ErrMissingCode                       = apperror.Validation("code is required")
	ErrMissingCodeVerifier               = apperror.Validation("code verifier is required")
	ErrInvalidGrantID                    = apperror.Validation("invalid grant id")
	ErrInvalidClientID                   = apperror.Validation("invalid client id")
	ErrInvalidCodeHash                   = apperror.Validation("invalid code hash")
	ErrInvalidRedirectURI                = apperror.Validation("invalid redirect uri")
	ErrInvalidCodeChallenge              = apperror.Validation("invalid code challenge")
	ErrInvalidResource                   = apperror.Validation("invalid resource")
	ErrInvalidScopes                     = apperror.Validation("invalid scopes")
)
