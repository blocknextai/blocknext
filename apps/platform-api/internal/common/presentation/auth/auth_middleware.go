package auth

import (
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/rbac"
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/account/application/sessions"
	commonApplicationAuth "github.com/blocknextai/platform-api/internal/common/application/auth"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	authTokenPrefix          = "Bearer "
	webSocketAuthTokenPrefix = "bearer."
)

var (
	ErrAuthorizationHeaderRequired = apperror.Unauthorized("authorization header required")
	ErrInvalidToken                = apperror.Unauthorized("invalid token")
	ErrSessionRevoked              = apperror.Unauthorized("session revoked")
	ErrNotAuthenticated            = apperror.Unauthorized("not authenticated")
	ErrPermissionRequired          = apperror.Forbidden("permission required")
	ErrInvalidOrganizationID       = apperror.Validation("invalid organization id")
	ErrSessionCheck                = apperror.Internal("session check failed")
	ErrPermissionCheck             = apperror.Internal("permission check failed")
)

type AuthMiddleware struct {
	authJWTService                jwt.AuthJWTService
	userPermissionChecker         commonApplicationAuth.UserPermissionChecker
	organizationPermissionChecker commonApplicationAuth.OrganizationPermissionChecker
	sessionService                accountApplicationSessions.SessionService
}

func NewAuthMiddleware(
	authJWTService jwt.AuthJWTService,
	userPermissionChecker commonApplicationAuth.UserPermissionChecker,
	organizationPermissionChecker commonApplicationAuth.OrganizationPermissionChecker,
	sessionService accountApplicationSessions.SessionService,
) *AuthMiddleware {
	return &AuthMiddleware{
		authJWTService:                authJWTService,
		userPermissionChecker:         userPermissionChecker,
		organizationPermissionChecker: organizationPermissionChecker,
		sessionService:                sessionService,
	}
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString, err := m.getToken(c)
		if err != nil {
			return err
		}
		return m.authenticate(c, tokenString)
	}
}

func (m *AuthMiddleware) AuthenticateWebSocket() fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString := tokenFromWebSocketProtocol(c.Get("Sec-WebSocket-Protocol"))
		if tokenString == "" {
			tokenString = strings.TrimSpace(c.Query("token"))
		}
		if tokenString == "" {
			return ErrAuthorizationHeaderRequired
		}
		return m.authenticate(c, tokenString)
	}
}

func (m *AuthMiddleware) RequireUserPermission(requiredPermission *rbac.Permission) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)
		if userID == uuid.Nil {
			return ErrNotAuthenticated
		}

		hasPermission, err := m.userPermissionChecker.HasPermission(c.RequestCtx(), userID, requiredPermission)
		if err != nil {
			return ErrPermissionCheck.WithCause(err)
		}
		if !hasPermission {
			return ErrPermissionRequired
		}

		return c.Next()
	}
}

func (m *AuthMiddleware) RequireOrganizationPermission(requiredPermission *rbac.Permission) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)
		if userID == uuid.Nil {
			return ErrNotAuthenticated
		}

		organizationID, err := uuid.Parse(c.Params("organizationId"))
		if err != nil {
			return ErrInvalidOrganizationID
		}

		hasPermission, err := m.organizationPermissionChecker.HasPermission(c.RequestCtx(), organizationID, userID, requiredPermission)
		if err != nil {
			return ErrPermissionCheck.WithCause(err)
		}
		if !hasPermission {
			return ErrPermissionRequired
		}

		commonHTTP.SetOrganizationID(c, organizationID)

		return c.Next()
	}
}

func (m *AuthMiddleware) authenticate(c fiber.Ctx, tokenString string) error {
	claims, err := m.authJWTService.ValidateAccessToken(tokenString)
	if err != nil {
		return ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ErrInvalidToken
	}

	commonHTTP.SetUserID(c, userID)

	if claims.SessionID != "" {
		if sessionID, err := uuid.Parse(claims.SessionID); err == nil {
			revoked, err := m.sessionService.IsSessionRevoked(c.RequestCtx(), sessionID)
			if err != nil {
				return ErrSessionCheck.WithCause(err)
			}
			if revoked {
				return ErrSessionRevoked
			}
			commonHTTP.SetSessionID(c, sessionID)
		}
	}

	return c.Next()
}

func (m *AuthMiddleware) getToken(c fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if strings.TrimSpace(authHeader) == "" {
		return "", ErrAuthorizationHeaderRequired
	}

	if !strings.HasPrefix(authHeader, authTokenPrefix) {
		return "", ErrInvalidToken
	}

	tokenString := strings.TrimPrefix(authHeader, authTokenPrefix)
	if strings.TrimSpace(tokenString) == "" {
		return "", ErrInvalidToken
	}

	return tokenString, nil
}

func tokenFromWebSocketProtocol(header string) string {
	if header == "" {
		return ""
	}
	for raw := range strings.SplitSeq(header, ",") {
		if token, ok := strings.CutPrefix(strings.TrimSpace(raw), webSocketAuthTokenPrefix); ok {
			return token
		}
	}
	return ""
}
