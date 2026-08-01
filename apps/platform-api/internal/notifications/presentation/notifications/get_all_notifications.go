package notifications

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/getallnotifications"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetAllUserNotificationsRequest struct {
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllUserNotificationsHandler(handler cqrs.Handler[*getallnotifications.GetAllNotificationsQuery, *getallnotifications.GetAllNotificationsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllUserNotificationsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getallnotifications.GetAllNotificationsQuery{
			UserID:         userID,
			OrganizationID: nil,
			Search:         searchRequest,
			Pagination:     paginationRequest,
		})
		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}

type GetAllOrganizationNotificationsRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	resultPkg.SearchRequest
	resultPkg.PaginationRequest
}

func NewGetAllOrganizationNotificationsHandler(handler cqrs.Handler[*getallnotifications.GetAllNotificationsQuery, *getallnotifications.GetAllNotificationsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetAllOrganizationNotificationsRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		searchRequest := request.SearchRequest.Normalize()
		paginationRequest := request.PaginationRequest.Normalize()

		result, err := handler.Handle(c.RequestCtx(), &getallnotifications.GetAllNotificationsQuery{
			UserID:         userID,
			OrganizationID: new(request.OrganizationID),
			Search:         searchRequest,
			Pagination:     paginationRequest,
		})
		if err != nil {
			return err
		}

		return commonHTTP.RespondPaginated(c, result.Items, result.TotalCount, paginationRequest)
	}
}
