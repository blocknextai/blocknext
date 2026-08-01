package taskexecutions

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/gettaskexecutionbyid"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetTaskExecutionByIDRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	ExecutionID    uuid.UUID `uri:"executionId"`
}

func NewGetExecutionByIDHandler(handler cqrs.Handler[*gettaskexecutionbyid.GetTaskExecutionByIDQuery, *gettaskexecutionbyid.GetTaskExecutionByIDResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetTaskExecutionByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &gettaskexecutionbyid.GetTaskExecutionByIDQuery{
			ID:             request.ExecutionID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
