package credentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/getcredentialsfornodes"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetUserCredentialsForNodesRequest struct {
	NodeIDs []string `query:"nodeIds"`
}

func NewGetUserCredentialsForNodesHandler(handler cqrs.Handler[*getcredentialsfornodes.GetCredentialsForNodesQuery, *getcredentialsfornodes.GetCredentialsForNodesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetUserCredentialsForNodesRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &getcredentialsfornodes.GetCredentialsForNodesQuery{
			OwnerType: commonDomain.OwnerTypeUser,
			OwnerID:   userID,
			NodeIDs:   request.NodeIDs,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}

type GetOrganizationCredentialsForNodesRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	NodeIDs        []string  `query:"nodeIds"`
}

func NewGetOrganizationCredentialsForNodesHandler(handler cqrs.Handler[*getcredentialsfornodes.GetCredentialsForNodesQuery, *getcredentialsfornodes.GetCredentialsForNodesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationCredentialsForNodesRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &getcredentialsfornodes.GetCredentialsForNodesQuery{
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
