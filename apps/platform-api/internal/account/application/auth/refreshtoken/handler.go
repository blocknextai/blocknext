package refreshtoken

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/auth/jwt"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/account/application/sessions"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/google/uuid"
)

type Handler struct {
	userRepository accountDomainUsers.UserRepository
	sessionService accountApplicationSessions.SessionService
	authJWTService jwt.AuthJWTService
}

func New(
	userRepository accountDomainUsers.UserRepository,
	sessionService accountApplicationSessions.SessionService,
	authJWTService jwt.AuthJWTService,
) *Handler {
	return &Handler{
		userRepository: userRepository,
		sessionService: sessionService,
		authJWTService: authJWTService,
	}
}

func (h *Handler) Handle(ctx context.Context, command *RefreshTokenCommand) (*accountApplicationAuth.AccessTokenResponse, error) {
	claims, err := h.authJWTService.ValidateRefreshToken(command.RefreshToken)
	if err != nil {
		return nil, accountApplicationAuth.ErrInvalidRefreshToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, accountApplicationAuth.ErrInvalidRefreshToken
	}

	user, err := h.userRepository.GetByID(ctx, userID)
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

	session, err := h.sessionService.ValidateRefreshTokenForRotation(ctx, sessionID, command.RefreshToken)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := h.authJWTService.GenerateAccessToken(userID, sessionID)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToGenerateToken.WithCause(err)
	}

	newRefreshToken, err := h.authJWTService.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		return nil, accountApplicationAuth.ErrFailedToGenerateToken.WithCause(err)
	}

	newExpiresAt := time.Now().UTC().Add(h.authJWTService.RefreshTokenTTL())

	if err = h.sessionService.CommitRefreshTokenRotation(ctx, session, newRefreshToken, newExpiresAt); err != nil {
		return nil, err
	}

	return &accountApplicationAuth.AccessTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
