package passwordchanged

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/application/idempotency"
	notificationsApplication "github.com/blocknextai/platform-api/internal/notifications/application/notifications"
	notificationsDomainNotifications "github.com/blocknextai/platform-api/internal/notifications/domain/notifications"
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

func (h *Handler) Handle(ctx context.Context, event accountDomainUsers.PasswordChangedDomainEvent) error {
	return h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		return h.eventBusInboxService.EnsureOnce(txCtx, "notifications:account.user.password_changed", func(txCtx context.Context) error {
			return h.notificationService.Create(txCtx, notificationsApplication.CreateInput{
				Type:         "account.password_changed",
				Level:        notificationsDomainNotifications.LevelWarning,
				AudienceType: notificationsDomainNotifications.AudienceTypeUser,
				AudienceID:   event.UserID,
				Title:        "Your password was changed",
				Data:         map[string]any{},
			})
		})
	})
}
