package taskexecutions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	taskexecutionsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases/taskexecutions"
)

type GetTaskExecutionByIDRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	ExecutionID    uuid.UUID `uri:"executionId"`
}

func NewGetExecutionByIDHandler(service *taskexecutionsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetTaskExecutionByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetTaskExecutionByID(c.RequestCtx(), &taskexecutionsUseCases.GetTaskExecutionByIDQuery{
			ID:             request.ExecutionID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
