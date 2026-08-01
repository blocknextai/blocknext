package organizationusercreated

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/rbac"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/application/idempotency"
	notificationsApplication "github.com/blocknextai/platform-api/internal/notifications/application/notifications"
	notificationsDomainNotifications "github.com/blocknextai/platform-api/internal/notifications/domain/notifications"
	organizationsDomainOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
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
		eventBus:             eventBus,
		eventBusInboxService: eventBusInboxService,
		notificationService:  notificationService,
		transactionManager:   transactionManager,
	}
	eventbus.Subscribe(eventBus, handler.Handle)
	return handler
}

func (h *Handler) Handle(ctx context.Context, event organizationsDomainOrganizationUsers.OrganizationUserCreatedDomainEvent) error {
	if event.Role == rbac.OrganizationOwnerRole.Name {
		return nil
	}

	return h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		return h.eventBusInboxService.EnsureOnce(txCtx, "notifications:organizations.organization_user.created", func(txCtx context.Context) error {
			return h.notificationService.Create(txCtx, notificationsApplication.CreateInput{
				Type:         "organization.member_added",
				Level:        notificationsDomainNotifications.LevelInfo,
				AudienceType: notificationsDomainNotifications.AudienceTypeUser,
				AudienceID:   event.UserID,
				Title:        "You've been added to an organization",
				Data: map[string]any{
					"organizationId": event.OrganizationID,
					"role":           event.Role,
				},
			})
		})
	})
}
