package publishing

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/json"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	outboxMessagesDomain "github.com/blocknextai/platform-api/internal/eventbus/domain/outboxmessages"
)

type PublisherService interface {
	Enqueue(ctx context.Context, event commonDomain.DomainEvent) error
}

type publisherService struct {
	repository outboxMessagesDomain.Repository
}

func NewPublisherService(repository outboxMessagesDomain.Repository) PublisherService {
	return &publisherService{
		repository: repository,
	}
}

func (s *publisherService) Enqueue(ctx context.Context, event commonDomain.DomainEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return ErrMarshalEvent
	}

	message, err := outboxMessagesDomain.New(event.EventName(), payload, time.Now().UTC())
	if err != nil {
		return err
	}

	if err := s.repository.Create(ctx, message); err != nil {
		return err
	}

	return nil
}
