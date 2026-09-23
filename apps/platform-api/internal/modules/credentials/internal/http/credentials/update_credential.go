package credentials

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/credentials/internal/usecases/credentials"
)

type UpdateOrganizationCredentialRequest struct {
	OrganizationID uuid.UUID      `uri:"organizationId"`
	CredentialID   uuid.UUID      `uri:"credentialId"`
	Key            string         `json:"key"`
	Title          string         `json:"title"`
	Data           map[string]any `json:"data"`
}

func NewUpdateOrganizationCredentialHandler(service *credentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateOrganizationCredentialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.UpdateCredential(c.RequestCtx(), &credentialsUseCases.UpdateCredentialCommand{
			ID:        request.CredentialID,
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			Key:       request.Key,
			Title:     request.Title,
			Data:      request.Data,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("credential updated")))
	}
}
