package magiclink

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/magiclink/consume"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type MagicLinkConsumeRequest struct {
	Token string `json:"token"`
}

func NewMagicLinkConsumeHandler(handler cqrs.Handler[*consume.MagicLinkConsumeCommand, *accountApplicationAuth.AccessTokenResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		req := new(MagicLinkConsumeRequest)
		if err := c.Bind().All(req); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &consume.MagicLinkConsumeCommand{
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
