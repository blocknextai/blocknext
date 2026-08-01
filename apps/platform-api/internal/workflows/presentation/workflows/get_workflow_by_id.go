package workflows

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows/getworkflowbyid"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetWorkflowByIDRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	WorkflowID     uuid.UUID `uri:"workflowId"`
}

func NewGetWorkflowByIDHandler(handler cqrs.Handler[*getworkflowbyid.GetWorkflowByIDQuery, *getworkflowbyid.GetWorkflowByIDResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetWorkflowByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &getworkflowbyid.GetWorkflowByIDQuery{
			OrganizationID: request.OrganizationID,
			WorkflowID:     request.WorkflowID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
