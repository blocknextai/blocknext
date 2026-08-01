package organizations

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/getallorganizations"
	"github.com/gofiber/fiber/v3"
)

func NewGetAllOrganizationsHandler(handler cqrs.Handler[*getallorganizations.GetAllOrganizationsQuery, *getallorganizations.GetAllOrganizationsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &getallorganizations.GetAllOrganizationsQuery{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
