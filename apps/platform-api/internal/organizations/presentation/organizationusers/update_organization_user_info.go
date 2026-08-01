package organizationusers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/updateorganizationuserinfo"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateOrganizationUserInfoRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	UserID         uuid.UUID `uri:"userId"`
	Alias          *string   `json:"alias"`
}

func NewUpdateOrganizationUserInfoHandler(handler cqrs.Handler[*updateorganizationuserinfo.UpdateOrganizationUserInfoCommand, *updateorganizationuserinfo.UpdateOrganizationUserInfoResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateOrganizationUserInfoRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &updateorganizationuserinfo.UpdateOrganizationUserInfoCommand{
			UserID:             userID,
			OrganizationID:     request.OrganizationID,
			OrganizationUserID: request.UserID,
			Alias:              request.Alias,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("organization user info updated")))
	}
}
