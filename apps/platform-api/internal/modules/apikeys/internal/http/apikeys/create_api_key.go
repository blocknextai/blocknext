package apikeys

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
	apikeysUseCases "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/usecases/apikeys"
)

type CreateOrganizationAPIKeyRequest struct {
	OrganizationID uuid.UUID                   `uri:"organizationId"`
	Name           string                      `json:"name"`
	Scopes         apiKeysDomainAPIKeys.Scopes `json:"scopes"`
}

func NewCreateOrganizationAPIKeyHandler(service *apikeysUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateOrganizationAPIKeyRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.CreateAPIKey(c.RequestCtx(), &apikeysUseCases.CreateAPIKeyCommand{
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			Name:      request.Name,
			Scopes:    request.Scopes,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("api key created")))
	}
}
