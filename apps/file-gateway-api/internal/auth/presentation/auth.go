package presentation

import (
	authDomainToken "github.com/blocknextai/file-gateway-api/internal/auth/domain/token"
	"github.com/blocknextai/go-packages/result"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, tokenService authDomainToken.Service) {
	router.Post("/auth/token", func(c fiber.Ctx) error {
		token, err := tokenService.Generate()
		if err != nil {
			return err
		}

		return c.JSON(result.Ok(AuthTokenResponse{
			Token: token,
		}))
	})
}
