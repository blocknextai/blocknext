package sessions

import (
	"time"

	"github.com/blocknextai/go-packages/database"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/google/uuid"
)

type Session struct {
	database.BaseEntity

	UserID                uuid.UUID
	AuthProvider          accountDomain.AuthProvider
	IPAddress             string
	UserAgent             string
	IsRevoked             bool
	RefreshTokenHash      *string
	RefreshTokenExpiresAt *time.Time
	TokenFamily           *uuid.UUID
	TokenGeneration       int
}

func New(
	id uuid.UUID,
	userID uuid.UUID,
	authProvider accountDomain.AuthProvider,
	ipAddress string,
	userAgent string,
	refreshTokenHash string,
	refreshTokenExpiresAt time.Time,
	tokenFamily uuid.UUID,
) (*Session, error) {
	utcNow := time.Now().UTC()

	session := &Session{
		BaseEntity: database.BaseEntity{
			ID:        id,
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		UserID:                userID,
		AuthProvider:          authProvider,
		IPAddress:             ipAddress,
		UserAgent:             userAgent,
		IsRevoked:             false,
		RefreshTokenHash:      new(refreshTokenHash),
		RefreshTokenExpiresAt: new(refreshTokenExpiresAt),
		TokenFamily:           new(tokenFamily),
		TokenGeneration:       0,
	}

	return session.validateThenReturn()
}

func (s *Session) Revoke() (*Session, error) {
	s.UpdatedAt = time.Now().UTC()

	s.IsRevoked = true

	return s.validateThenReturn()
}

func (s *Session) RotateRefreshToken(newHash string, expiresAt time.Time) (rotated *Session, previousGeneration int, err error) {
	previousGeneration = s.TokenGeneration

	s.UpdatedAt = time.Now().UTC()
	s.RefreshTokenHash = new(newHash)
	s.RefreshTokenExpiresAt = new(expiresAt)
	s.TokenGeneration++

	rotated, err = s.validateThenReturn()
	return rotated, previousGeneration, err
}

func (s *Session) validateThenReturn() (*Session, error) {
	return s, nil
}
