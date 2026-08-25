package toolinvocations

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/executions/application/toolinvocations/gettoolinvocationbyid"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetToolInvocationByIDRequest struct {
	OrganizationID   uuid.UUID `uri:"organizationId"`
	ToolInvocationID uuid.UUID `uri:"toolInvocationId"`
}

func NewGetToolInvocationByIDHandler(handler cqrs.Handler[*gettoolinvocationbyid.GetToolInvocationByIDQuery, *gettoolinvocationbyid.GetToolInvocationByIDResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetToolInvocationByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &gettoolinvocationbyid.GetToolInvocationByIDQuery{
			ID:             request.ToolInvocationID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
