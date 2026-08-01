package taskrunner

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/taskrunner/application/canceltask"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CancelTaskRequest struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewCancelTaskHandler(handler *canceltask.CancelTaskHandler) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CancelTaskRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &canceltask.CancelTaskCommand{
			TriggeredByUserID: userID,
			OrganizationID:    request.OrganizationID,
			ID:                request.ID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("task cancelled")))
	}
}
