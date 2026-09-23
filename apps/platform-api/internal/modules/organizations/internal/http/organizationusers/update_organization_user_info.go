package organizationusers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	organizationusersUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases/organizationusers"
)

type UpdateOrganizationUserInfoRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	UserID         uuid.UUID `uri:"userId"`
	Alias          *string   `json:"alias"`
}

func NewUpdateOrganizationUserInfoHandler(service *organizationusersUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateOrganizationUserInfoRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.UpdateOrganizationUserInfo(c.RequestCtx(), &organizationusersUseCases.UpdateOrganizationUserInfoCommand{
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
