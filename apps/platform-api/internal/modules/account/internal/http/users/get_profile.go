package users

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	usersUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/users"
)

func NewGetProfileHandler(service *usersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := service.GetProfile(c.RequestCtx(), &usersUseCases.GetProfileQuery{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
