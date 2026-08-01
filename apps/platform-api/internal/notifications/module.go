package notifications

import (
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/application/idempotency"
	notificationsApplication "github.com/blocknextai/platform-api/internal/notifications/application/notifications"
	notificationsInfrastructure "github.com/blocknextai/platform-api/internal/notifications/infrastructure"
	notificationsInfrastructureRecipients "github.com/blocknextai/platform-api/internal/notifications/infrastructure/notificationrecipients"
	notificationsInfrastructureNotifications "github.com/blocknextai/platform-api/internal/notifications/infrastructure/notifications"
	notificationsPresentation "github.com/blocknextai/platform-api/internal/notifications/presentation"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                   *sql.DB
	TransactionManager   database.TransactionManager
	EventBus             *eventbus.Bus
	EventBusInboxService *idempotency.InboxService

	OrganizationUserService organizationsApplicationOrganizationUsers.OrganizationUserService
}

type Module struct {
	handlers *notificationsInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	notificationRepository := notificationsInfrastructureNotifications.NewNotificationRepository(deps.DB)
	notificationRecipientRepository := notificationsInfrastructureRecipients.NewNotificationRecipientRepository(deps.DB)

	notificationService := notificationsApplication.NewNotificationService(
		notificationRepository,
		notificationRecipientRepository,
		deps.OrganizationUserService,
	)

	handlers := notificationsInfrastructure.RegisterInfrastructure(notificationsInfrastructure.RegisterInfrastructureDeps{
		TransactionManager:   deps.TransactionManager,
		EventBus:             deps.EventBus,
		EventBusInboxService: deps.EventBusInboxService,

		NotificationRecipientRepository: notificationRecipientRepository,
		NotificationService:             notificationService,
	})

	return &Module{
		handlers: handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	notificationsPresentation.RegisterPresentation(router, authMiddleware, m.handlers)
}
