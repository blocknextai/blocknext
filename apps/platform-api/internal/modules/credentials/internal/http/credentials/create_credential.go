package credentials

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/credentials/internal/usecases/credentials"
)

type CreateOrganizationCredentialRequest struct {
	OrganizationID uuid.UUID                               `uri:"organizationId"`
	SourceType     credentialsDomainCredentials.SourceType `json:"sourceType"`
	Key            string                                  `json:"key"`
	Title          string                                  `json:"title"`
	Data           map[string]any                          `json:"data"`
}

func NewCreateOrganizationCredentialHandler(service *credentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateOrganizationCredentialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		sourceType := request.SourceType
		if sourceType == "" {
			sourceType = credentialsDomainCredentials.SourceTypeOwner
		}

		result, err := service.CreateCredential(c.RequestCtx(), &credentialsUseCases.CreateCredentialCommand{
			OwnerType:  commonDomain.OwnerTypeOrganization,
			OwnerID:    request.OrganizationID,
			SourceType: sourceType,
			Key:        request.Key,
			Title:      request.Title,
			Data:       request.Data,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("credential created")))
	}
}
