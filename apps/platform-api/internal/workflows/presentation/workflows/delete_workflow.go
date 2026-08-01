package workflows

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/deleteworkflow"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteWorkflowRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	WorkflowID     uuid.UUID `uri:"workflowId"`
}

func NewDeleteWorkflowHandler(handler cqrs.Handler[*deleteworkflow.DeleteWorkflowCommand, *deleteworkflow.DeleteWorkflowResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteWorkflowRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &deleteworkflow.DeleteWorkflowCommand{
			OrganizationID: request.OrganizationID,
			WorkflowID:     request.WorkflowID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("workflow deleted")))
	}
}
