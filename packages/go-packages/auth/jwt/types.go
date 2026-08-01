// Package jwt issues and validates HS256-signed access and refresh JWTs that
// carry a user and session identity, with configurable issuer, audience,
// expiration, and clock leeway.
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType distinguishes between access and refresh tokens.
type TokenType string

const (
	// TokenTypeAccessToken identifies an access token.
	TokenTypeAccessToken TokenType = "access"
	// TokenTypeRefreshToken identifies a refresh token.
	TokenTypeRefreshToken TokenType = "refresh"
)

// AccessTokenClaims holds the registered and custom claims carried by an
// access token.
type AccessTokenClaims struct {
	jwt.RegisteredClaims
	TokenType TokenType `json:"tokenType"`
	SessionID string    `json:"sessionId"`
}

// RefreshTokenClaims holds the registered and custom claims carried by a
// refresh token.
type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	TokenType TokenType `json:"tokenType"`
	SessionID string    `json:"sessionId"`
}

// AuthJWTService generates and validates access and refresh tokens and
// reports their configured lifetimes.
type AuthJWTService interface {
	// GenerateAccessToken returns a signed access token for the given user and
	// session.
	GenerateAccessToken(userID uuid.UUID, sessionID uuid.UUID) (string, error)
	// GenerateRefreshToken returns a signed refresh token for the given user
	// and session.
	GenerateRefreshToken(userID uuid.UUID, sessionID uuid.UUID) (string, error)

	// ValidateAccessToken parses and verifies tokenString and returns its
	// access token claims.
	ValidateAccessToken(tokenString string) (*AccessTokenClaims, error)
	// ValidateRefreshToken parses and verifies tokenString and returns its
	// refresh token claims.
	ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error)

	// AccessTokenTTL returns the configured access token lifetime.
	AccessTokenTTL() time.Duration
	// RefreshTokenTTL returns the configured refresh token lifetime.
	RefreshTokenTTL() time.Duration
}
