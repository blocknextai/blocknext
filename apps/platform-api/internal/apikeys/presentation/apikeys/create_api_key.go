package apikeys

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/createapikey"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
)

type CreateOrganizationAPIKeyRequest struct {
	OrganizationID uuid.UUID                   `uri:"organizationId"`
	Name           string                      `json:"name"`
	Scopes         apiKeysDomainAPIKeys.Scopes `json:"scopes"`
}

func NewCreateOrganizationAPIKeyHandler(handler cqrs.Handler[*createapikey.CreateAPIKeyCommand, *createapikey.CreateAPIKeyResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateOrganizationAPIKeyRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &createapikey.CreateAPIKeyCommand{
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
