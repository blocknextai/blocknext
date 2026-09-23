package toolinvocations

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	toolinvocationsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases/toolinvocations"
)

type GetToolInvocationByIDRequest struct {
	OrganizationID   uuid.UUID `uri:"organizationId"`
	ToolInvocationID uuid.UUID `uri:"toolInvocationId"`
}

func NewGetToolInvocationByIDHandler(service *toolinvocationsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetToolInvocationByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetToolInvocationByID(c.RequestCtx(), &toolinvocationsUseCases.GetToolInvocationByIDQuery{
			ID:             request.ToolInvocationID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
