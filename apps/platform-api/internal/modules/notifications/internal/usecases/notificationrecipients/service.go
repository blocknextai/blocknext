package notificationrecipients

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/modules/notifications/internal/domain/notificationrecipients"
)

type Service struct {
	notificationRecipientRepository notificationrecipients.NotificationRecipientRepository
	transactionManager              database.TransactionManager
}

func NewService(
	notificationRecipientRepository notificationrecipients.NotificationRecipientRepository,
	transactionManager database.TransactionManager,
) *Service {
	return &Service{
		notificationRecipientRepository: notificationRecipientRepository,
		transactionManager:              transactionManager,
	}
}
