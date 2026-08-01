package organizations

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/getorganizationbyid"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetOrganizationByIDRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewGetOrganizationByIDHandler(handler cqrs.Handler[*getorganizationbyid.GetOrganizationByIDQuery, *getorganizationbyid.GetOrganizationByIDResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &getorganizationbyid.GetOrganizationByIDQuery{
			UserID:         userID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
