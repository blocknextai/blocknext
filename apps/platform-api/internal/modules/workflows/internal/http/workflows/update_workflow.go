package workflows

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	workflowsUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/workflows"
)

type UpdateWorkflowRequest struct {
	OrganizationID uuid.UUID  `uri:"organizationId"`
	WorkflowID     uuid.UUID  `uri:"workflowId"`
	Title          *string    `json:"title"`
	Description    *string    `json:"description"`
	IsPinned       *bool      `json:"isPinned"`
	Nodes          []dag.Node `json:"nodes"`
	Edges          []dag.Edge `json:"edges"`
}

func NewUpdateWorkflowHandler(service *workflowsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateWorkflowRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.UpdateWorkflow(c.RequestCtx(), &workflowsUseCases.UpdateWorkflowCommand{
			OrganizationID: request.OrganizationID,
			WorkflowID:     request.WorkflowID,
			Title:          request.Title,
			Description:    request.Description,
			IsPinned:       request.IsPinned,
			Nodes:          request.Nodes,
			Edges:          request.Edges,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("workflow updated")))
	}
}
