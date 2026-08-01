package organizationusers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/deleteorganizationuser"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteOrganizationUserRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	UserID         uuid.UUID `uri:"userId"`
}

func NewDeleteOrganizationUserHandler(handler cqrs.Handler[*deleteorganizationuser.DeleteOrganizationUserCommand, *deleteorganizationuser.DeleteOrganizationUserResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteOrganizationUserRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &deleteorganizationuser.DeleteOrganizationUserCommand{
			OrganizationID: request.OrganizationID,
			UserID:         request.UserID,
			ForceDelete:    false,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("organization user deleted")))
	}
}
