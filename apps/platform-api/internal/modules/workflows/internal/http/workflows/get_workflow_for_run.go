package workflows

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	workflowsUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/workflows"
)

type GetWorkflowForRunRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	WorkflowID     uuid.UUID `uri:"workflowId"`
}

func NewGetWorkflowForRunHandler(service *workflowsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetWorkflowForRunRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetWorkflowForRun(c.RequestCtx(), &workflowsUseCases.GetWorkflowForRunQuery{
			OrganizationID: request.OrganizationID,
			WorkflowID:     request.WorkflowID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
