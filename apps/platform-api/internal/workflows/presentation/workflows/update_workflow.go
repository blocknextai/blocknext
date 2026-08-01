package workflows

import (
	"github.com/blocknextai/go-packages/dag"
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/updateworkflow"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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

func NewUpdateWorkflowHandler(handler cqrs.Handler[*updateworkflow.UpdateWorkflowCommand, *updateworkflow.UpdateWorkflowResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateWorkflowRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &updateworkflow.UpdateWorkflowCommand{
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
