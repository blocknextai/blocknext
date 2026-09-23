package http

import (
	"github.com/gofiber/fiber/v3"

	authToken "github.com/blocknextai/file-gateway-api/internal/modules/auth/internal/token"
	"github.com/blocknextai/go-packages/result"
)

func RegisterRoutes(router fiber.Router, tokenService authToken.Service) {
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
