package taskrunner

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/taskrunner/application/rerunfailed"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RerunFailedRequest struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewRerunFailedHandler(handler *rerunfailed.RerunFailedHandler) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RerunFailedRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &rerunfailed.RerunFailedCommand{
			TriggeredByUserID: userID,
			OrganizationID:    request.OrganizationID,
			ID:                request.ID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("task rerun started")))
	}
}
