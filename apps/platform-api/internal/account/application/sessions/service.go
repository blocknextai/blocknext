package sessions

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/digest"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainSessions "github.com/blocknextai/platform-api/internal/account/domain/sessions"
	"github.com/google/uuid"
)

const (
	sessionStatusKeyPrefix = "session_status:"
	sessionStatusActive    = "active"
	sessionStatusRevoked   = "revoked"
)

type SessionService interface {
	TrackSession(
		ctx context.Context,
		sessionID uuid.UUID,
		userID uuid.UUID,
		authProvider accountDomain.AuthProvider,
		ipAddress string,
		userAgent string,
		refreshToken string,
		refreshTokenExpiresAt time.Time,
		tokenFamily uuid.UUID,
	) error

	IsSessionRevoked(ctx context.Context, sessionID uuid.UUID) (bool, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	RevokeTokenFamily(ctx context.Context, tokenFamily uuid.UUID) ([]uuid.UUID, error)
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	ValidateRefreshTokenForRotation(ctx context.Context, sessionID uuid.UUID, oldToken string) (*accountDomainSessions.Session, error)
	CommitRefreshTokenRotation(ctx context.Context, session *accountDomainSessions.Session, newToken string, newExpiresAt time.Time) error
}

type sessionService struct {
	sessionRepository accountDomainSessions.SessionRepository
	cacheService      cache.Service
	tokenTTL          time.Duration
}

func NewSessionService(sessionRepository accountDomainSessions.SessionRepository, cacheService cache.Service, tokenTTL time.Duration) SessionService {
	return &sessionService{
		sessionRepository: sessionRepository,
		cacheService:      cacheService,
		tokenTTL:          tokenTTL,
	}
}

func (s *sessionService) TrackSession(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
	authProvider accountDomain.AuthProvider,
	ipAddress string,
	userAgent string,
	refreshToken string,
	refreshTokenExpiresAt time.Time,
	tokenFamily uuid.UUID,
) error {
	refreshTokenHash := digest.SHA256Hex(refreshToken)

	session, err := accountDomainSessions.New(sessionID, userID, authProvider, ipAddress, userAgent, refreshTokenHash, refreshTokenExpiresAt, tokenFamily)
	if err != nil {
		return err
	}

	return s.sessionRepository.Create(ctx, session)
}

func (s *sessionService) sessionStatusKey(sessionID uuid.UUID) string {
	sessionIDStr := sessionID.String()

	var builder strings.Builder
	builder.Grow(len(sessionStatusKeyPrefix) + len(sessionIDStr))
	builder.WriteString(sessionStatusKeyPrefix)
	builder.WriteString(sessionIDStr)

	return builder.String()
}

func (s *sessionService) IsSessionRevoked(ctx context.Context, sessionID uuid.UUID) (bool, error) {
	key := s.sessionStatusKey(sessionID)

	status, err := s.cacheService.Get(ctx, key)
	if err == nil && status != "" {
		return status == sessionStatusRevoked, nil
	}

	session, err := s.sessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return false, err
	}

	if session.IsRevoked {
		_ = s.cacheService.Set(ctx, key, sessionStatusRevoked, s.tokenTTL)
		return true, nil
	}

	_ = s.cacheService.Set(ctx, key, sessionStatusActive, s.tokenTTL)
	return false, nil
}

func (s *sessionService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	if err := s.sessionRepository.RevokeByID(ctx, sessionID, time.Now().UTC()); err != nil {
		return err
	}

	_ = s.cacheService.Set(ctx, s.sessionStatusKey(sessionID), sessionStatusRevoked, s.tokenTTL)

	return nil
}

func (s *sessionService) RevokeTokenFamily(ctx context.Context, tokenFamily uuid.UUID) ([]uuid.UUID, error) {
	ids, err := s.sessionRepository.RevokeAllByTokenFamily(ctx, tokenFamily, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		_ = s.cacheService.Set(ctx, s.sessionStatusKey(id), sessionStatusRevoked, s.tokenTTL)
	}

	return ids, nil
}

func (s *sessionService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	ids, err := s.sessionRepository.RevokeAllByUserID(ctx, userID, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		_ = s.cacheService.Set(ctx, s.sessionStatusKey(id), sessionStatusRevoked, s.tokenTTL)
	}

	return ids, nil
}

func (s *sessionService) ValidateRefreshTokenForRotation(ctx context.Context, sessionID uuid.UUID, oldToken string) (*accountDomainSessions.Session, error) {
	oldTokenHash := digest.SHA256Hex(oldToken)

	session, err := s.sessionRepository.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.IsRevoked {
		return nil, accountDomainSessions.ErrSessionAlreadyRevoked
	}

	if session.RefreshTokenHash == nil || *session.RefreshTokenHash != oldTokenHash {
		if session.TokenFamily != nil {
			_, _ = s.RevokeTokenFamily(ctx, *session.TokenFamily)
		} else {
			_ = s.RevokeSession(ctx, sessionID)
		}
		return nil, accountDomainSessions.ErrRefreshTokenReused
	}

	if session.RefreshTokenExpiresAt != nil && time.Now().UTC().After(*session.RefreshTokenExpiresAt) {
		return nil, accountDomainSessions.ErrRefreshTokenExpired
	}

	return session, nil
}

func (s *sessionService) CommitRefreshTokenRotation(ctx context.Context, session *accountDomainSessions.Session, newToken string, newExpiresAt time.Time) error {
	newTokenHash := digest.SHA256Hex(newToken)

	rotated, previousGeneration, err := session.RotateRefreshToken(newTokenHash, newExpiresAt)
	if err != nil {
		return err
	}

	err = s.sessionRepository.UpdateRefreshToken(ctx, rotated, previousGeneration)
	if errors.Is(err, accountDomainSessions.ErrRefreshTokenReused) {
		if rotated.TokenFamily != nil {
			_, _ = s.RevokeTokenFamily(ctx, *rotated.TokenFamily)
		} else {
			_ = s.RevokeSession(ctx, rotated.ID)
		}
		return accountDomainSessions.ErrRefreshTokenReused
	}
	return err
}
