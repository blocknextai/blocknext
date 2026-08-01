package tokenissuer

import (
	"context"
	"log/slog"
	"time"

	"github.com/blocknextai/go-packages/auth/jwt"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/account/application/sessions"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/google/uuid"
)

type Service interface {
	IssueTokens(
		ctx context.Context,
		userID uuid.UUID,
		authProvider accountDomain.AuthProvider,
		ipAddress string,
		userAgent string,
	) (*accountApplicationAuth.AccessTokenResponse, error)
}

type service struct {
	authJWTService jwt.AuthJWTService
	sessionService accountApplicationSessions.SessionService
	userRepository accountDomainUsers.UserRepository
}

func NewService(
	authJWTService jwt.AuthJWTService,
	sessionService accountApplicationSessions.SessionService,
	userRepository accountDomainUsers.UserRepository,
) Service {
	return &service{
		authJWTService: authJWTService,
		sessionService: sessionService,
		userRepository: userRepository,
	}
}

func (s *service) IssueTokens(
	ctx context.Context,
	userID uuid.UUID,
	authProvider accountDomain.AuthProvider,
	ipAddress string,
	userAgent string,
) (*accountApplicationAuth.AccessTokenResponse, error) {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.IsBanned {
		return nil, accountApplicationAuth.ErrUserBanned
	}

	sessionID := bnuuid.NewV7()

	accessToken, err := s.authJWTService.GenerateAccessToken(userID, sessionID)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToGenerateToken
	}

	refreshToken, err := s.authJWTService.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToGenerateToken
	}

	tokenFamily := bnuuid.NewV7()
	refreshTokenExpiresAt := time.Now().UTC().Add(s.authJWTService.RefreshTokenTTL())

	if trackErr := s.sessionService.TrackSession(ctx, sessionID, userID, authProvider, ipAddress, userAgent, refreshToken, refreshTokenExpiresAt, tokenFamily); trackErr != nil {
		slog.ErrorContext(ctx, "Failed to track session",
			"component", "TokenIssuer",
			"user_id", userID,
			"session_id", sessionID,
			"error", trackErr)
	}

	return &accountApplicationAuth.AccessTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
