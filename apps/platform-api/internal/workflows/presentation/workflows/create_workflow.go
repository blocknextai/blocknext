package workflows

import (
	"github.com/blocknextai/go-packages/dag"
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/createworkflow"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateWorkflowRequest struct {
	OrganizationID uuid.UUID  `uri:"organizationId"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Nodes          []dag.Node `json:"nodes"`
	Edges          []dag.Edge `json:"edges"`
}

func NewCreateWorkflowHandler(handler cqrs.Handler[*createworkflow.CreateWorkflowCommand, *createworkflow.CreateWorkflowResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateWorkflowRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &createworkflow.CreateWorkflowCommand{
			OrganizationID: request.OrganizationID,
			UserID:         userID,
			Title:          request.Title,
			Description:    request.Description,
			Nodes:          request.Nodes,
			Edges:          request.Edges,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("workflow created")))
	}
}
