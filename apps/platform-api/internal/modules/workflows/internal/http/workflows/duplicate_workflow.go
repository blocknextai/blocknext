package workflows

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	workflowsUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/workflows"
)

type DuplicateWorkflowRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	WorkflowID     uuid.UUID `uri:"workflowId"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
}

func NewDuplicateWorkflowHandler(service *workflowsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DuplicateWorkflowRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.DuplicateWorkflow(c.RequestCtx(), &workflowsUseCases.DuplicateWorkflowCommand{
			WorkflowID:     request.WorkflowID,
			OrganizationID: request.OrganizationID,
			UserID:         userID,
			Title:          request.Title,
			Description:    request.Description,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("workflow duplicated")))
	}
}
