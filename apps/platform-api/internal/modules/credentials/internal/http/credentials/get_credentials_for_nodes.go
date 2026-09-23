package credentials

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/credentials/internal/usecases/credentials"
)

type GetOrganizationCredentialsForNodesRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	NodeIDs        []string  `query:"nodeIds"`
}

func NewGetOrganizationCredentialsForNodesHandler(service *credentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationCredentialsForNodesRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetCredentialsForNodes(c.RequestCtx(), &credentialsUseCases.GetCredentialsForNodesQuery{
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			NodeIDs:   request.NodeIDs,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
