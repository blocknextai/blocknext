package workflows

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/duplicateworkflow"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DuplicateWorkflowRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	WorkflowID     uuid.UUID `uri:"workflowId"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
}

func NewDuplicateWorkflowHandler(handler cqrs.Handler[*duplicateworkflow.DuplicateWorkflowCommand, *duplicateworkflow.DuplicateWorkflowResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DuplicateWorkflowRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &duplicateworkflow.DuplicateWorkflowCommand{
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
