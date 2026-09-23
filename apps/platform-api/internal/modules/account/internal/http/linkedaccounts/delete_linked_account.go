package linkedaccounts

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	linkedaccountsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/linkedaccounts"
)

type DeleteLinkedAccountRequest struct {
	LinkedAccountID uuid.UUID `uri:"linkedAccountId"`
}

func NewDeleteLinkedAccountHandler(service *linkedaccountsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)
		request := new(DeleteLinkedAccountRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.DeleteLinkedAccount(c.RequestCtx(), &linkedaccountsUseCases.DeleteLinkedAccountCommand{
			UserID:          userID,
			LinkedAccountID: request.LinkedAccountID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("linked account deleted")))
	}
}
