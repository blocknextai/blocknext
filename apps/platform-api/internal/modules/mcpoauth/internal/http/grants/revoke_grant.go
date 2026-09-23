package grants

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	grantsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/grants"
)

type RevokeUserGrantRequest struct {
	GrantID uuid.UUID `uri:"grantId"`
}

func NewRevokeUserGrantHandler(service *grantsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RevokeUserGrantRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.RevokeGrant(c.RequestCtx(), &grantsUseCases.RevokeGrantCommand{
			GrantID: request.GrantID,
			UserID:  commonHTTP.GetUserID(c),
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("grant revoked")))
	}
}
