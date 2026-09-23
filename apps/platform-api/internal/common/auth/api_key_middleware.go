package auth

import (
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/apperror"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
)

const (
	APIKeyHeader    = "X-API-Key"
	APIKeyIDHeader  = "X-Auth-API-Key-ID"
	OwnerTypeHeader = "X-Auth-Owner-Type"
	OwnerIDHeader   = "X-Auth-Owner-ID"
	UserIDHeader    = "X-Auth-User-ID"

	apiKeyScopesContextKey = "apiKeyScopes"
)

var (
	ErrInvalidAPIKey       = apperror.Unauthorized("missing or invalid api key")
	ErrAPIKeyScopeRequired = apperror.Forbidden("the api key does not have the scope required to perform this operation")
)

type APIKeyMiddleware[S ~string] struct {
	validator APIKeyValidator[S]
}

func NewAPIKeyMiddleware[S ~string](
	validator APIKeyValidator[S],
) *APIKeyMiddleware[S] {
	return &APIKeyMiddleware[S]{
		validator: validator,
	}
}

func (m *APIKeyMiddleware[S]) Authenticate() fiber.Handler {
	return func(c fiber.Ctx) error {
		rawKey := strings.TrimSpace(c.Get(APIKeyHeader))
		if rawKey == "" {
			return ErrInvalidAPIKey
		}

		authenticated, err := m.validator.Validate(c.RequestCtx(), rawKey)
		if err != nil {
			return ErrInvalidAPIKey.WithCause(err)
		}

		c.Request().Header.Set(APIKeyIDHeader, authenticated.ID.String())
		c.Request().Header.Set(OwnerTypeHeader, authenticated.OwnerType.String())
		c.Request().Header.Set(OwnerIDHeader, authenticated.OwnerID.String())

		switch authenticated.OwnerType {
		case commonDomain.OwnerTypeUser:
			commonHTTP.SetUserID(c, authenticated.OwnerID)
		case commonDomain.OwnerTypeOrganization:
			commonHTTP.SetOrganizationID(c, authenticated.OwnerID)
		default:
			return ErrInvalidAPIKey
		}

		c.Locals(apiKeyScopesContextKey, authenticated.Scopes)

		return c.Next()
	}
}

func (m *APIKeyMiddleware[S]) RequireScope(scope S) fiber.Handler {
	return func(c fiber.Ctx) error {
		if scopes, ok := c.Locals(apiKeyScopesContextKey).([]S); ok && slices.Contains(scopes, scope) {
			return c.Next()
		}

		return ErrAPIKeyScopeRequired
	}
}
