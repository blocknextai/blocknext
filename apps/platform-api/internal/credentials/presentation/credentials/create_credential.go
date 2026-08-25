package credentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/createcredential"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateOrganizationCredentialRequest struct {
	OrganizationID uuid.UUID                               `uri:"organizationId"`
	SourceType     credentialsDomainCredentials.SourceType `json:"sourceType"`
	Key            string                                  `json:"key"`
	Title          string                                  `json:"title"`
	Data           map[string]any                          `json:"data"`
}

func NewCreateOrganizationCredentialHandler(handler cqrs.Handler[*createcredential.CreateCredentialCommand, *createcredential.CreateCredentialResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateOrganizationCredentialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		sourceType := request.SourceType
		if sourceType == "" {
			sourceType = credentialsDomainCredentials.SourceTypeOwner
		}

		result, err := handler.Handle(c.RequestCtx(), &createcredential.CreateCredentialCommand{
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
