package auth

import (
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	apiKeysApplicationAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/application/apikeys"
	apiKeysDomain "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

const (
	APIKeyHeader    = "X-API-Key"
	APIKeyIDHeader  = "X-Auth-API-Key-ID"
	OwnerTypeHeader = "X-Auth-Owner-Type"
	OwnerIDHeader   = "X-Auth-Owner-ID"

	apiKeyScopesContextKey = "apiKeyScopes"
)

var (
	ErrInvalidAPIKey       = apperror.Unauthorized("missing or invalid api key")
	ErrAPIKeyScopeRequired = apperror.Forbidden("the api key does not have the scope required to perform this operation")
)

type APIKeyMiddleware struct {
	validator apiKeysApplicationAPIKeys.APIKeyValidator
}

func NewAPIKeyMiddleware(
	validator apiKeysApplicationAPIKeys.APIKeyValidator,
) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		validator: validator,
	}
}

func (m *APIKeyMiddleware) Authenticate() fiber.Handler {
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

func (m *APIKeyMiddleware) RequireScope(scope apiKeysDomain.Scope) fiber.Handler {
	return func(c fiber.Ctx) error {
		if value := c.Locals(apiKeyScopesContextKey); value != nil {
			if scopes, ok := value.(apiKeysDomain.Scopes); ok {
				if scopes.Has(scope) {
					return c.Next()
				}
			}
		}

		return ErrAPIKeyScopeRequired
	}
}
