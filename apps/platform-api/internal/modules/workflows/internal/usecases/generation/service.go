package generation

import (
	"github.com/blocknextai/go-packages/database"
	generationDomainMessages "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/generation/messages"
	generationDomainSessions "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/generation/sessions"
)

type Service struct {
	sessionRepository  generationDomainSessions.SessionRepository
	messageRepository  generationDomainMessages.MessageRepository
	transactionManager database.TransactionManager
}

func NewService(
	sessionRepository generationDomainSessions.SessionRepository,
	messageRepository generationDomainMessages.MessageRepository,
	transactionManager database.TransactionManager,
) *Service {
	return &Service{
		sessionRepository:  sessionRepository,
		messageRepository:  messageRepository,
		transactionManager: transactionManager,
	}
}
