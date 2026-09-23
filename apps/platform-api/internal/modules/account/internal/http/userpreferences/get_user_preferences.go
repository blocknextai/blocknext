package userpreferences

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	userpreferencesUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/userpreferences"
)

func NewGetUserPreferencesHandler(service *userpreferencesUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := service.GetUserPreferences(c.RequestCtx(), &userpreferencesUseCases.GetUserPreferencesQuery{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
