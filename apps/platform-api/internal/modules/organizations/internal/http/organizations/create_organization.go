package organizations

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	organizationsUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizations"
)

type CreateOrganizationRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

func NewCreateOrganizationHandler(service *organizationsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateOrganizationRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.CreateOrganization(c.RequestCtx(), &organizationsUseCases.CreateOrganizationCommand{
			UserID:      userID,
			Title:       request.Title,
			Description: request.Description,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("organization created")))
	}
}
