package taskexecutions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	taskexecutionsUseCases "github.com/blocknextai/platform-api/internal/modules/executions/internal/usecases/taskexecutions"
)

type BulkDeleteTaskExecutionsRequest struct {
	OrganizationID uuid.UUID   `uri:"organizationId"`
	IDs            []uuid.UUID `json:"ids"`
}

func NewBulkDeleteTaskExecutionsHandler(service *taskexecutionsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(BulkDeleteTaskExecutionsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.BulkDeleteTaskExecutions(c.RequestCtx(), &taskexecutionsUseCases.BulkDeleteTaskExecutionsCommand{
			IDs:            request.IDs,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("task executions deleted")))
	}
}
