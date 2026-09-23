package magiclink

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	authUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/auth"
)

type MagicLinkConsumeRequest struct {
	Token string `json:"token"`
}

func NewMagicLinkConsumeHandler(service *authUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		req := new(MagicLinkConsumeRequest)
		if err := c.Bind().All(req); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.MagicLinkConsume(c.RequestCtx(), &authUseCases.MagicLinkConsumeCommand{
			Token:     req.Token,
			IPAddress: c.IP(),
			UserAgent: c.Get(fiber.HeaderUserAgent),
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("logged in")))
	}
}
