package jwt

import (
	"errors"
)

var (
	// ErrInvalidTokenType indicates the token's type does not match the
	// expected access or refresh type.
	ErrInvalidTokenType = errors.New("invalid token type")
	// ErrInvalidToken indicates the access token is malformed or failed
	// validation.
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken indicates the token has expired.
	ErrExpiredToken = errors.New("expired token")
	// ErrWeakSecretKey indicates the provided secret key is empty or shorter
	// than 32 bytes.
	ErrWeakSecretKey = errors.New("weak secret key")
	// ErrEmptyIssuer indicates the configured issuer is empty; a non-empty
	// issuer is required for issuer validation to be enforced.
	ErrEmptyIssuer = errors.New("empty issuer")
	// ErrEmptyAudience indicates the configured audience is empty; a non-empty
	// audience is required for audience validation to be enforced.
	ErrEmptyAudience = errors.New("empty audience")
	// ErrInvalidRefreshToken indicates the refresh token is malformed or
	// failed validation.
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// errUnexpectedSigningMethod indicates the token's signing method is not
	// the expected HMAC method. It is mapped to ErrInvalidToken /
	// ErrInvalidRefreshToken before reaching callers.
	errUnexpectedSigningMethod = errors.New("unexpected signing method")
)
