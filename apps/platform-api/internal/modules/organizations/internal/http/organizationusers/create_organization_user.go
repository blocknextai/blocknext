package organizationusers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	organizationusersUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizationusers"
)

type CreateOrganizationUserRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	Identifier     string    `json:"identifier"`
	Alias          string    `json:"alias"`
}

func NewCreateOrganizationUserHandler(service *organizationusersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateOrganizationUserRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.CreateOrganizationUser(c.RequestCtx(), &organizationusersUseCases.CreateOrganizationUserCommand{
			OrganizationID: request.OrganizationID,
			Identifier:     request.Identifier,
			Alias:          request.Alias,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("organization user added")))
	}
}
