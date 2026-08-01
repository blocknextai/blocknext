package eventbus

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/blocknextai/go-packages/json"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
)

type handler func(ctx context.Context, event commonDomain.DomainEvent) error

type decoder func(payload []byte) (commonDomain.DomainEvent, error)

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]handler
	decoders map[string]decoder
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]handler),
		decoders: make(map[string]decoder),
	}
}

func Subscribe[T commonDomain.DomainEvent](bus *Bus, handle func(ctx context.Context, event T) error) {
	var zero T
	name := zero.EventName()

	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.handlers[name] = append(bus.handlers[name], func(ctx context.Context, event commonDomain.DomainEvent) error {
		return handle(ctx, event.(T))
	})

	if _, ok := bus.decoders[name]; !ok {
		bus.decoders[name] = func(payload []byte) (commonDomain.DomainEvent, error) {
			var event T
			if err := json.Unmarshal(payload, &event); err != nil {
				return nil, err
			}
			return event, nil
		}
	}
}

func (bus *Bus) Dispatch(ctx context.Context, eventName string, payload []byte) error {
	bus.mu.RLock()
	decode := bus.decoders[eventName]
	handlers := bus.handlers[eventName]
	bus.mu.RUnlock()

	if decode == nil {
		slog.WarnContext(ctx, "no handler registered for event, skipping",
			"component", "eventbus",
			"event_name", eventName)
		return nil
	}

	event, err := decode(payload)
	if err != nil {
		return err
	}

	return fanOut(ctx, handlers, event)
}

func fanOut(ctx context.Context, handlers []handler, event commonDomain.DomainEvent) error {
	var errs []error
	for _, handle := range handlers {
		if err := handle(ctx, event); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
