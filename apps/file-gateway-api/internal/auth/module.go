package auth

import (
	authDomainToken "github.com/blocknextai/file-gateway-api/internal/auth/domain/token"
	authInfrastructureToken "github.com/blocknextai/file-gateway-api/internal/auth/infrastructure/token"
	authPresentation "github.com/blocknextai/file-gateway-api/internal/auth/presentation"
	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	ServiceKey string
	JWTService jwt.AuthJWTService
}

type Module struct {
	TokenService    authDomainToken.Service
	Middleware      fiber.Handler
	MatchServiceKey func(fiber.Ctx) bool
}

func NewModule(deps Dependencies) *Module {
	tokenService := authInfrastructureToken.NewService(deps.JWTService)
	middleware := authPresentation.NewAuthMiddleware(deps.ServiceKey, deps.JWTService)
	matchServiceKey := authPresentation.MatchServiceKey(deps.ServiceKey)

	return &Module{
		TokenService:    tokenService,
		Middleware:      middleware,
		MatchServiceKey: matchServiceKey,
	}
}

func (m *Module) Register(router fiber.Router) {
	authPresentation.RegisterRoutes(router, m.TokenService)
}
