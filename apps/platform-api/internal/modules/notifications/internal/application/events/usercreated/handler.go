package usercreated

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/idempotency"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	notificationsApplication "github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/notifications"
	notificationsDomainNotifications "github.com/blocknextai/platform-api/internal/modules/notifications/internal/domain/notifications"
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

func (h *Handler) Handle(ctx context.Context, event accountContract.UserCreatedDomainEvent) error {
	return h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		return h.eventBusInboxService.EnsureOnce(txCtx, "notifications:account.user.created", func(txCtx context.Context) error {
			return h.notificationService.Create(txCtx, notificationsApplication.CreateInput{
				Type:         "welcome",
				Level:        notificationsDomainNotifications.LevelSuccess,
				AudienceType: notificationsDomainNotifications.AudienceTypeUser,
				AudienceID:   event.UserID,
				Title:        "Welcome",
				Body:         new("Your account is ready. Start building your first workflow."),
			})
		})
	})
}
