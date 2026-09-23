package taskexecutions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	taskexecutionsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases/taskexecutions"
)

type DeleteTaskExecutionRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	ExecutionID    uuid.UUID `uri:"executionId"`
}

func NewDeleteExecutionHandler(service *taskexecutionsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteTaskExecutionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.DeleteTaskExecution(c.RequestCtx(), &taskexecutionsUseCases.DeleteTaskExecutionCommand{
			ID:             request.ExecutionID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("task execution deleted")))
	}
}
