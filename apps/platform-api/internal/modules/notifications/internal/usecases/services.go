package usecases

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/idempotency"
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/events/emailchanged"
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/events/organizationusercreated"
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/events/organizationuserrolechanged"
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/events/passwordchanged"
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/events/usercreated"
	notificationsApplication "github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/notifications"
	notificationsDomainRecipients "github.com/blocknextai/platform-api/internal/modules/notifications/internal/domain/notificationrecipients"
	notificationrecipientsUseCases "github.com/blocknextai/platform-api/internal/modules/notifications/internal/usecases/notificationrecipients"
)

type Services struct {
	Notificationrecipients *notificationrecipientsUseCases.Service
}

type ServiceDependencies struct {
	TransactionManager   database.TransactionManager
	EventBus             *eventbus.Bus
	EventBusInboxService *idempotency.InboxService

	NotificationRecipientRepository notificationsDomainRecipients.NotificationRecipientRepository
	NotificationService             notificationsApplication.NotificationService
}

func NewServices(deps ServiceDependencies) *Services {
	usercreated.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)
	organizationusercreated.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)
	passwordchanged.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)
	emailchanged.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)
	organizationuserrolechanged.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)

	return &Services{
		Notificationrecipients: notificationrecipientsUseCases.NewService(deps.NotificationRecipientRepository, deps.TransactionManager),
	}
}
