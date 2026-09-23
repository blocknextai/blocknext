package auth

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth"
)

type RefreshTokenCommand struct {
	RefreshToken string
	IPAddress    string
	UserAgent    string
}

var (
	ErrRefreshTokenRequired = apperror.Validation("refresh token is required")
)

func (c *RefreshTokenCommand) Validate() error {
	if strings.TrimSpace(c.RefreshToken) == "" {
		return ErrRefreshTokenRequired
	}

	return nil
}

func (s *Service) RefreshToken(ctx context.Context, command *RefreshTokenCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	claims, err := s.authJWTService.ValidateRefreshToken(command.RefreshToken)
	if err != nil {
		return nil, accountApplicationAuth.ErrInvalidRefreshToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, accountApplicationAuth.ErrInvalidRefreshToken
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.IsBanned {
		return nil, accountApplicationAuth.ErrUserBanned
	}

	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return nil, accountApplicationAuth.ErrInvalidRefreshToken
	}

	session, err := s.sessionService.ValidateRefreshTokenForRotation(ctx, sessionID, command.RefreshToken)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := s.authJWTService.GenerateAccessToken(userID, sessionID)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToGenerateToken.WithCause(err)
	}

	newRefreshToken, err := s.authJWTService.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToGenerateToken.WithCause(err)
	}

	newExpiresAt := time.Now().UTC().Add(s.authJWTService.RefreshTokenTTL())

	if err = s.sessionService.CommitRefreshTokenRotation(ctx, session, newRefreshToken, newExpiresAt); err != nil {
		return nil, err
	}

	return &accountApplicationAuth.AccessTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
