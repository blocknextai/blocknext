package notifications

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/idempotency"
	notificationsApplication "github.com/blocknextai/platform-api/internal/modules/notifications/internal/application/notifications"
	notificationsHTTP "github.com/blocknextai/platform-api/internal/modules/notifications/internal/http"
	notificationsPostgresRecipients "github.com/blocknextai/platform-api/internal/modules/notifications/internal/postgres/notificationrecipients"
	notificationsPostgresNotifications "github.com/blocknextai/platform-api/internal/modules/notifications/internal/postgres/notifications"
	notificationsUseCases "github.com/blocknextai/platform-api/internal/modules/notifications/internal/usecases"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
)

type Dependencies struct {
	DB                   *sql.DB
	TransactionManager   database.TransactionManager
	EventBus             *eventbus.Bus
	EventBusInboxService *idempotency.InboxService

	OrganizationUserService organizationsContract.OrganizationUserService
}

type Module struct {
	useCases *notificationsUseCases.Services
}

func NewModule(deps Dependencies) *Module {
	notificationRepository := notificationsPostgresNotifications.NewNotificationRepository(deps.DB)
	notificationRecipientRepository := notificationsPostgresRecipients.NewNotificationRecipientRepository(deps.DB)

	notificationService := notificationsApplication.NewNotificationService(
		notificationRepository,
		notificationRecipientRepository,
		deps.OrganizationUserService,
	)

	useCases := notificationsUseCases.NewServices(notificationsUseCases.ServiceDependencies{
		TransactionManager:   deps.TransactionManager,
		EventBus:             deps.EventBus,
		EventBusInboxService: deps.EventBusInboxService,

		NotificationRecipientRepository: notificationRecipientRepository,
		NotificationService:             notificationService,
	})

	return &Module{
		useCases: useCases,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	notificationsHTTP.RegisterRoutes(router, authMiddleware, m.useCases)
}
