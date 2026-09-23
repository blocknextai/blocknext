package linkedaccounts

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	linkedaccountsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/linkedaccounts"
)

func NewGetAllLinkedAccountsHandler(service *linkedaccountsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := service.GetAllLinkedAccounts(c.RequestCtx(), &linkedaccountsUseCases.GetAllLinkedAccountsQuery{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
