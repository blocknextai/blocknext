package usernonces

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrUserNonceNotFound             = apperror.NotFound("user nonce not found")
	ErrNonceIsRequired               = apperror.Validation("nonce is required")
	ErrCodeVerifierIsRequired        = apperror.Validation("code verifier is required")
	ErrCodeChallengeIsRequired       = apperror.Validation("code challenge is required")
	ErrCodeChallengeMethodIsRequired = apperror.Validation("code challenge method is required")
	ErrExpiresAtIsRequired           = apperror.Validation("expires at is required")
	ErrExpiresAtInThePast            = apperror.Validation("expires at is in the past")
	ErrNonceAlreadyExists            = apperror.Conflict("nonce already exists")
)
