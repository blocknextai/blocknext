package linkedaccounts

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/linkedaccounts/deletelinkedaccount"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteLinkedAccountRequest struct {
	LinkedAccountID uuid.UUID `uri:"linkedAccountId"`
}

func NewDeleteLinkedAccountHandler(handler cqrs.Handler[*deletelinkedaccount.DeleteLinkedAccountCommand, *deletelinkedaccount.DeleteLinkedAccountResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)
		request := new(DeleteLinkedAccountRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &deletelinkedaccount.DeleteLinkedAccountCommand{
			UserID:          userID,
			LinkedAccountID: request.LinkedAccountID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("linked account deleted")))
	}
}
