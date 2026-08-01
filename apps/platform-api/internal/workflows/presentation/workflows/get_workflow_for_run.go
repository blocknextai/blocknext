package workflows

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/getworkflowforrun"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetWorkflowForRunRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	WorkflowID     uuid.UUID `uri:"workflowId"`
}

func NewGetWorkflowForRunHandler(handler cqrs.Handler[*getworkflowforrun.GetWorkflowForRunQuery, *getworkflowforrun.GetWorkflowForRunResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetWorkflowForRunRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &getworkflowforrun.GetWorkflowForRunQuery{
			OrganizationID: request.OrganizationID,
			WorkflowID:     request.WorkflowID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
