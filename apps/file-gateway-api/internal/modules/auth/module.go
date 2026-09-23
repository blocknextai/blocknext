package auth

import (
	"github.com/gofiber/fiber/v3"

	authHTTP "github.com/blocknextai/file-gateway-api/internal/modules/auth/internal/http"
	authToken "github.com/blocknextai/file-gateway-api/internal/modules/auth/internal/token"
	"github.com/blocknextai/go-packages/auth/jwt"
)

type Dependencies struct {
	ServiceKey string
	JWTService jwt.AuthJWTService
}

type Module struct {
	TokenService    authToken.Service
	Middleware      fiber.Handler
	MatchServiceKey func(fiber.Ctx) bool
}

func NewModule(deps Dependencies) *Module {
	tokenService := authToken.NewService(deps.JWTService)
	middleware := authHTTP.NewAuthMiddleware(deps.ServiceKey, deps.JWTService)
	matchServiceKey := authHTTP.MatchServiceKey(deps.ServiceKey)

	return &Module{
		TokenService:    tokenService,
		Middleware:      middleware,
		MatchServiceKey: matchServiceKey,
	}
}

func (m *Module) Register(router fiber.Router) {
	authHTTP.RegisterRoutes(router, m.TokenService)
}
