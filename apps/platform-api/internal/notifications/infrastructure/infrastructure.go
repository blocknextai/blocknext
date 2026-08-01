package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/eventbus/application/idempotency"
	"github.com/blocknextai/platform-api/internal/notifications/application/events/emailchanged"
	"github.com/blocknextai/platform-api/internal/notifications/application/events/organizationusercreated"
	"github.com/blocknextai/platform-api/internal/notifications/application/events/organizationuserrolechanged"
	"github.com/blocknextai/platform-api/internal/notifications/application/events/passwordchanged"
	"github.com/blocknextai/platform-api/internal/notifications/application/events/usercreated"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/deletenotification"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/getallnotifications"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/getnotificationcounts"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/markallnotificationsread"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/markallnotificationsseen"
	"github.com/blocknextai/platform-api/internal/notifications/application/notificationrecipients/marknotificationread"
	notificationsApplication "github.com/blocknextai/platform-api/internal/notifications/application/notifications"
	notificationsDomainRecipients "github.com/blocknextai/platform-api/internal/notifications/domain/notificationrecipients"
)

type Handlers struct {
	GetAllNotifications      cqrs.Handler[*getallnotifications.GetAllNotificationsQuery, *getallnotifications.GetAllNotificationsResponse]
	GetNotificationCounts    cqrs.Handler[*getnotificationcounts.GetNotificationCountsQuery, *getnotificationcounts.GetNotificationCountsResponse]
	MarkNotificationRead     cqrs.Handler[*marknotificationread.MarkNotificationReadCommand, *marknotificationread.MarkNotificationReadResponse]
	MarkAllNotificationsRead cqrs.Handler[*markallnotificationsread.MarkAllNotificationsReadCommand, *markallnotificationsread.MarkAllNotificationsReadResponse]
	MarkAllNotificationsSeen cqrs.Handler[*markallnotificationsseen.MarkAllNotificationsSeenCommand, *markallnotificationsseen.MarkAllNotificationsSeenResponse]
	DeleteNotification       cqrs.Handler[*deletenotification.DeleteNotificationCommand, *deletenotification.DeleteNotificationResponse]
}

type RegisterInfrastructureDeps struct {
	TransactionManager   database.TransactionManager
	EventBus             *eventbus.Bus
	EventBusInboxService *idempotency.InboxService

	NotificationRecipientRepository notificationsDomainRecipients.NotificationRecipientRepository
	NotificationService             notificationsApplication.NotificationService
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	usercreated.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)
	organizationusercreated.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)
	passwordchanged.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)
	emailchanged.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)
	organizationuserrolechanged.New(deps.NotificationService, deps.EventBus, deps.EventBusInboxService, deps.TransactionManager)

	return &Handlers{
		GetAllNotifications:      cqrs.ValidationBehavior(getallnotifications.New(deps.NotificationRecipientRepository)),
		GetNotificationCounts:    cqrs.ValidationBehavior(getnotificationcounts.New(deps.NotificationRecipientRepository)),
		MarkNotificationRead:     cqrs.ValidationBehavior(marknotificationread.New(deps.NotificationRecipientRepository, deps.TransactionManager)),
		MarkAllNotificationsRead: cqrs.ValidationBehavior(markallnotificationsread.New(deps.NotificationRecipientRepository, deps.TransactionManager)),
		MarkAllNotificationsSeen: cqrs.ValidationBehavior(markallnotificationsseen.New(deps.NotificationRecipientRepository, deps.TransactionManager)),
		DeleteNotification:       cqrs.ValidationBehavior(deletenotification.New(deps.NotificationRecipientRepository, deps.TransactionManager)),
	}
}
