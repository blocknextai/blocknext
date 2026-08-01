package taskexecutions

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/bulkdeletetaskexecutions"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type BulkDeleteTaskExecutionsRequest struct {
	OrganizationID uuid.UUID   `uri:"organizationId"`
	IDs            []uuid.UUID `json:"ids"`
}

func NewBulkDeleteTaskExecutionsHandler(handler cqrs.Handler[*bulkdeletetaskexecutions.BulkDeleteTaskExecutionsCommand, *bulkdeletetaskexecutions.BulkDeleteTaskExecutionsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(BulkDeleteTaskExecutionsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &bulkdeletetaskexecutions.BulkDeleteTaskExecutionsCommand{
			IDs:            request.IDs,
			OrganizationID: request.OrganizationID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("task executions deleted")))
	}
}
