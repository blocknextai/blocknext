package credentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/getallcredentials"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetAllUserCredentialsRequest struct {
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllUserCredentialsHandler(handler cqrs.Handler[*getallcredentials.GetAllCredentialsQuery, *getallcredentials.GetAllCredentialsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllUserCredentialsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getallcredentials.GetAllCredentialsQuery{
			OwnerType:  commonDomain.OwnerTypeUser,
			OwnerID:    userID,
			Search:     searchRequest,
			Pagination: paginationRequest,
		})
		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}

type GetAllOrganizationCredentialsRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllOrganizationCredentialsHandler(handler cqrs.Handler[*getallcredentials.GetAllCredentialsQuery, *getallcredentials.GetAllCredentialsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllOrganizationCredentialsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getallcredentials.GetAllCredentialsQuery{
			OwnerType:  commonDomain.OwnerTypeOrganization,
			OwnerID:    request.OrganizationID,
			Search:     searchRequest,
			Pagination: paginationRequest,
		})
		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}
