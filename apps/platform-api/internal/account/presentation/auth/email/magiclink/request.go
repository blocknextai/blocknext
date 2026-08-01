package magiclink

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/magiclink/request"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type MagicLinkRequestRequest struct {
	Email string `json:"email"`
}

func NewMagicLinkRequestHandler(handler cqrs.Handler[*request.MagicLinkRequestCommand, *request.MagicLinkRequestResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		req := new(MagicLinkRequestRequest)
		if err := c.Bind().All(req); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &request.MagicLinkRequestCommand{
			Email: req.Email,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("magic link sent")))
	}
}
