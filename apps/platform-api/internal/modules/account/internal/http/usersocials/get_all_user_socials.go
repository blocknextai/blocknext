package usersocials

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	usersocialsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/usersocials"
)

func NewGetAllUserSocialsHandler(service *usersocialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := service.GetAllUserSocials(c.RequestCtx(), &usersocialsUseCases.GetAllUserSocialsQuery{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
