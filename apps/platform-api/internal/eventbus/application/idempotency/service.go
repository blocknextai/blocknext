package idempotency

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/database"
	inboxEntriesDomain "github.com/blocknextai/platform-api/internal/eventbus/domain/inboxentries"
)

type InboxService struct {
	repository         inboxEntriesDomain.Repository
	transactionManager database.TransactionManager
}

func NewInboxService(repository inboxEntriesDomain.Repository, transactionManager database.TransactionManager) *InboxService {
	return &InboxService{
		repository:         repository,
		transactionManager: transactionManager,
	}
}

func (s *InboxService) EnsureOnce(ctx context.Context, handlerKey string, fn func(ctx context.Context) error) error {
	eventID, ok := EventIDFromContext(ctx)
	if !ok {
		return fn(ctx)
	}

	entry, err := inboxEntriesDomain.New(handlerKey, eventID, time.Now().UTC())
	if err != nil {
		return err
	}

	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		inserted, err := s.repository.MarkProcessed(txCtx, entry)
		if err != nil {
			return err
		}
		if !inserted {
			return nil
		}

		return fn(txCtx)
	})
}
