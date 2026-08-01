package taskrunner

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/taskrunner/application/rerunall"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RerunAllRequest struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewRerunAllHandler(handler *rerunall.RerunAllHandler) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RerunAllRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &rerunall.RerunAllCommand{
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
