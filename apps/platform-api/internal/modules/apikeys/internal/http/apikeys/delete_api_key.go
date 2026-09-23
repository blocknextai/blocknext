package apikeys

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	apikeysUseCases "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/usecases/apikeys"
)

type DeleteOrganizationAPIKeyRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	APIKeyID       uuid.UUID `uri:"apiKeyId"`
}

func NewDeleteOrganizationAPIKeyHandler(service *apikeysUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteOrganizationAPIKeyRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.DeleteAPIKey(c.RequestCtx(), &apikeysUseCases.DeleteAPIKeyCommand{
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			APIKeyID:  request.APIKeyID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("api key deleted")))
	}
}
