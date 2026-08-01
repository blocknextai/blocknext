package organizationusers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/createorganizationuser"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateOrganizationUserRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	Identifier     string    `json:"identifier"`
	Alias          string    `json:"alias"`
}

func NewCreateOrganizationUserHandler(handler cqrs.Handler[*createorganizationuser.CreateOrganizationUserCommand, *createorganizationuser.CreateOrganizationUserResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateOrganizationUserRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &createorganizationuser.CreateOrganizationUserCommand{
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
