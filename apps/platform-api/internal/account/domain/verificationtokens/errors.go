package verificationtokens

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrVerificationTokenNotFound      = apperror.NotFound("verification token not found")
	ErrVerificationTokenExpired       = apperror.Validation("verification token expired")
	ErrInvalidPurpose                 = apperror.Validation("invalid verification token purpose")
	ErrUserIDIsRequired               = apperror.Validation("user id is required")
	ErrTokenHashIsRequired            = apperror.Validation("token hash is required")
	ErrEmailIsRequired                = apperror.Validation("email is required")
	ErrVerificationTokenAlreadyExists = apperror.Conflict("verification token already exists")
)
