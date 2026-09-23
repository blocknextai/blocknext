package organizationuserrolechanged

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/idempotency"
	notificationsApplication "github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/notifications"
	notificationsDomainNotifications "github.com/blocknextai/platform-api/internal/modules/notifications/internal/domain/notifications"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
)

type Handler struct {
	notificationService  notificationsApplication.NotificationService
	eventBus             *eventbus.Bus
	eventBusInboxService *idempotency.InboxService
	transactionManager   database.TransactionManager
}

func New(
	notificationService notificationsApplication.NotificationService,
	eventBus *eventbus.Bus,
	eventBusInboxService *idempotency.InboxService,
	transactionManager database.TransactionManager,
) *Handler {
	handler := &Handler{
		notificationService:  notificationService,
		eventBus:             eventBus,
		eventBusInboxService: eventBusInboxService,
		transactionManager:   transactionManager,
	}
	eventbus.Subscribe(eventBus, handler.Handle)
	return handler
}

func (h *Handler) Handle(ctx context.Context, event organizationsContract.OrganizationUserRoleChangedDomainEvent) error {
	return h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		return h.eventBusInboxService.EnsureOnce(txCtx, "notifications:organizations.organization_user.role_changed", func(txCtx context.Context) error {
			return h.notificationService.Create(txCtx, notificationsApplication.CreateInput{
				Type:         "organization.role_changed",
				Level:        notificationsDomainNotifications.LevelInfo,
				AudienceType: notificationsDomainNotifications.AudienceTypeUser,
				AudienceID:   event.UserID,
				Title:        "Your role was changed",
				Data: map[string]any{
					"organizationId": event.OrganizationID,
					"newRole":        event.NewRole,
				},
			})
		})
	})
}
