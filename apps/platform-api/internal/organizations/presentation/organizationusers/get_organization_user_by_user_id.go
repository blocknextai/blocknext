package organizationusers

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizationusers/getorganizationuserbyuserid"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetOrganizationUserByUserIDRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	UserID         uuid.UUID `uri:"userId"`
}

func NewGetOrganizationUserByUserIDHandler(handler cqrs.Handler[*getorganizationuserbyuserid.GetOrganizationUserByUserIDQuery, *getorganizationuserbyuserid.GetOrganizationUserByUserIDResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationUserByUserIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &getorganizationuserbyuserid.GetOrganizationUserByUserIDQuery{
			OrganizationID: request.OrganizationID,
			UserID:         request.UserID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
