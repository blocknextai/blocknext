package taskrunner

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	tasksUseCases "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/usecases/tasks"
)

type RerunAllRequest struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `uri:"organizationId"`
}

func NewRerunAllHandler(service *tasksUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RerunAllRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.RerunAll(c.RequestCtx(), &tasksUseCases.RerunAllCommand{
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
