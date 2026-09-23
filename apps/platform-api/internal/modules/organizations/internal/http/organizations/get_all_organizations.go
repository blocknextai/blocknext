package organizations

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	organizationsUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizations"
)

func NewGetAllOrganizationsHandler(service *organizationsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := service.GetAllOrganizations(c.RequestCtx(), &organizationsUseCases.GetAllOrganizationsQuery{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
