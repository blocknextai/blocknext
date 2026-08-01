package presentation

import (
	"crypto/subtle"
	"strings"

	authDomain "github.com/blocknextai/file-gateway-api/internal/auth/domain"
	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/gofiber/fiber/v3"
)

func MatchServiceKey(serviceKey string) func(fiber.Ctx) bool {
	serviceKeyBytes := []byte(serviceKey)
	return func(c fiber.Ctx) bool {
		key := c.Get(HeaderServiceKey)
		if key == "" {
			return false
		}
		return subtle.ConstantTimeCompare([]byte(key), serviceKeyBytes) == 1
	}
}

func NewAuthMiddleware(serviceKey string, jwtService jwt.AuthJWTService) fiber.Handler {
	matchServiceKey := MatchServiceKey(serviceKey)

	return func(c fiber.Ctx) error {
		if c.Get(HeaderServiceKey) != "" {
			if matchServiceKey(c) {
				return c.Next()
			}
			return authDomain.ErrInvalidServiceKey
		}

		auth := c.Get(HeaderAuthorization)
		if after, ok := strings.CutPrefix(auth, BearerPrefix); ok {
			token := after
			if _, err := jwtService.ValidateAccessToken(token); err != nil {
				return authDomain.ErrInvalidToken.WithCause(err)
			}
			return c.Next()
		}

		return authDomain.ErrMissingToken
	}
}
