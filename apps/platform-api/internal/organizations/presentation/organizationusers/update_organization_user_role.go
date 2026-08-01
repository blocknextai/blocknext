package organizationusers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/updateorganizationuserrole"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateOrganizationUserRoleRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	UserID         uuid.UUID `uri:"userId"`
	Role           string    `json:"role"`
}

func NewUpdateOrganizationUserRoleHandler(handler cqrs.Handler[*updateorganizationuserrole.UpdateOrganizationUserRoleCommand, *updateorganizationuserrole.UpdateOrganizationUserRoleResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateOrganizationUserRoleRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &updateorganizationuserrole.UpdateOrganizationUserRoleCommand{
			UserID:             userID,
			OrganizationID:     request.OrganizationID,
			OrganizationUserID: request.UserID,
			Role:               request.Role,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("organization user role updated")))
	}
}
