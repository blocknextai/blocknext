package magiclink

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type MagicLinkRequestRequest struct {
	Email string `json:"email"`
}

func NewMagicLinkRequestHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		req := new(MagicLinkRequestRequest)
		if err := c.Bind().All(req); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.MagicLinkRequest(c.RequestCtx(), &authUseCases.MagicLinkRequestCommand{
			Email: req.Email,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("magic link sent")))
	}
}
