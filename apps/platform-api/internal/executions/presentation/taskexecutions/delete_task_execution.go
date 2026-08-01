package taskexecutions

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/deletetaskexecution"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteTaskExecutionRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	ExecutionID    uuid.UUID `uri:"executionId"`
}

func NewDeleteExecutionHandler(handler cqrs.Handler[*deletetaskexecution.DeleteTaskExecutionCommand, *deletetaskexecution.DeleteTaskExecutionResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteTaskExecutionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &deletetaskexecution.DeleteTaskExecutionCommand{
			ID:             request.ExecutionID,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("task execution deleted")))
	}
}
