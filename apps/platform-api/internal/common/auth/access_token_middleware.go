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
	ScopesHeader = "X-Auth-Scopes"

	bearerTokenPrefix           = "Bearer "
	accessTokenScopesContextKey = "accessTokenScopes"
	invalidTokenErrorCode       = "invalid_token"
	insufficientScopeErrorCode  = "insufficient_scope"
)

var (
	ErrAccessTokenRequired      = apperror.Unauthorized("a valid access token is required")
	ErrAccessTokenScopeRequired = apperror.Forbidden("the access token does not have the scope required to perform this operation")
)

type AccessTokenMiddleware struct {
	validator                    AccessTokenValidator
	protectedResourceMetadataURL string
}

func NewAccessTokenMiddleware(
	validator AccessTokenValidator,
	protectedResourceMetadataURL string,
) *AccessTokenMiddleware {
	return &AccessTokenMiddleware{
		validator:                    validator,
		protectedResourceMetadataURL: strings.TrimSuffix(protectedResourceMetadataURL, "/"),
	}
}

func (m *AccessTokenMiddleware) Authenticate() fiber.Handler {
	return func(c fiber.Ctx) error {
		rawToken := bearerToken(c)
		if strings.TrimSpace(rawToken) == "" {
			return m.challenge(c, ErrAccessTokenRequired, "", nil)
		}

		authenticated, err := m.validator.Validate(c.RequestCtx(), rawToken, c.Path())
		if err != nil {
			return m.challenge(c, ErrAccessTokenRequired.WithCause(err), invalidTokenErrorCode, nil)
		}

		c.Request().Header.Set(OwnerTypeHeader, commonDomain.OwnerTypeOrganization.String())
		c.Request().Header.Set(OwnerIDHeader, authenticated.OrganizationID.String())
		c.Request().Header.Set(UserIDHeader, authenticated.UserID.String())
		c.Request().Header.Set(ScopesHeader, strings.Join(authenticated.Scopes, " "))

		commonHTTP.SetUserID(c, authenticated.UserID)
		commonHTTP.SetOrganizationID(c, authenticated.OrganizationID)

		c.Locals(accessTokenScopesContextKey, authenticated.Scopes)

		return c.Next()
	}
}

func (m *AccessTokenMiddleware) RequireScope(scope string) fiber.Handler {
	return m.RequireAnyScope(scope)
}

func (m *AccessTokenMiddleware) RequireAnyScope(scopes ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		for _, scope := range scopes {
			if HasAccessTokenScope(c, scope) {
				return c.Next()
			}
		}

		return m.challenge(c, ErrAccessTokenScopeRequired, insufficientScopeErrorCode, scopes)
	}
}

func IsAccessTokenRequest(c fiber.Ctx) bool {
	return strings.TrimSpace(bearerToken(c)) != ""
}

func HasAccessTokenScope(c fiber.Ctx, scope string) bool {
	granted, ok := c.Locals(accessTokenScopesContextKey).([]string)
	return ok && slices.Contains(granted, scope)
}

func (m *AccessTokenMiddleware) challenge(c fiber.Ctx, err error, errorCode string, scopes []string) error {
	var builder strings.Builder
	builder.WriteString(`Bearer resource_metadata="`)
	builder.WriteString(m.protectedResourceMetadataURL)
	builder.WriteString(c.Path())
	builder.WriteString(`"`)

	if strings.TrimSpace(errorCode) != "" {
		builder.WriteString(`, error="`)
		builder.WriteString(errorCode)
		builder.WriteString(`"`)
	}

	if len(scopes) > 0 {
		builder.WriteString(`, scope="`)
		builder.WriteString(strings.Join(scopes, " "))
		builder.WriteString(`"`)
	}

	c.Set(fiber.HeaderWWWAuthenticate, builder.String())

	return err
}

func bearerToken(c fiber.Ctx) string {
	token, ok := strings.CutPrefix(c.Get(fiber.HeaderAuthorization), bearerTokenPrefix)
	if !ok {
		return ""
	}
	return strings.TrimSpace(token)
}
