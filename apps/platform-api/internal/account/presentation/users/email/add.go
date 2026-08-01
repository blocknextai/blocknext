package email

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/auth/email/add"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type AddEmailRequest struct {
	Email string `json:"email"`
}

func NewAddEmailHandler(handler cqrs.Handler[*add.AddEmailCommand, *add.AddEmailResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(AddEmailRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &add.AddEmailCommand{
			UserID: userID,
			Email:  request.Email,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusAccepted).JSON(resultPkg.Ok(result, resultPkg.WithMessage("email added")))
	}
}
