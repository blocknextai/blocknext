package email

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/resendverification"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type ResendVerificationRequest struct {
	Email string `json:"email"`
}

func NewResendVerificationHandler(handler cqrs.Handler[*resendverification.ResendVerificationCommand, *resendverification.ResendVerificationResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(ResendVerificationRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &resendverification.ResendVerificationCommand{
			Email: request.Email,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("verification email sent")))
	}
}
