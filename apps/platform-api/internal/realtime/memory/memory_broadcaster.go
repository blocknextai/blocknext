package memory

import (
	"context"
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"
	"sync"

	"github.com/blocknextai/platform-api/internal/realtime/events"
	taskRunnerDomainNode "github.com/blocknextai/platform-api/internal/taskrunner/domain/node"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
)

type memoryBroadcaster struct {
	mu          sync.RWMutex
	subscribers map[uuid.UUID]map[chan string]struct{}
	closed      bool
}

func New() *memoryBroadcaster {
	return &memoryBroadcaster{
		subscribers: make(map[uuid.UUID]map[chan string]struct{}),
	}
}

func (b *memoryBroadcaster) Ping(_ context.Context) error {
	return nil
}

func (b *memoryBroadcaster) PublishTaskEvent(_ context.Context, event *taskRunnerDomainTask.TaskEvent) error {
	payload, err := events.MarshalTask(event)
	if err != nil {
		return err
	}

	b.publish(event.OrganizationID, payload)

	return nil
}

func (b *memoryBroadcaster) PublishNodeEvent(_ context.Context, event *taskRunnerDomainNode.NodeEvent) error {
	payload, err := events.MarshalNode(event)
	if err != nil {
		return err
	}

	b.publish(event.OrganizationID, payload)

	return nil
}

func (b *memoryBroadcaster) PublishToolInvocationEvent(_ context.Context, event *executionsDomainToolInvocations.ToolInvocationEvent) error {
	payload, err := events.MarshalToolInvocation(event)
	if err != nil {
		return err
	}

	b.publish(event.OrganizationID, payload)

	return nil
}

func (b *memoryBroadcaster) publish(organizationID uuid.UUID, payload string) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for subscriber := range b.subscribers[organizationID] {
		select {
		case subscriber <- payload:
		default:
			// The subscriber is behind; dropping keeps the publisher moving.
		}
	}
}

func (b *memoryBroadcaster) Subscribe(ctx context.Context, organizationID uuid.UUID) (<-chan string, error) {
	subscriber := make(chan string, events.SubscriberBuffer)

	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		close(subscriber)
		return subscriber, nil
	}
	if b.subscribers[organizationID] == nil {
		b.subscribers[organizationID] = make(map[chan string]struct{})
	}
	b.subscribers[organizationID][subscriber] = struct{}{}
	b.mu.Unlock()

	go func() {
		<-ctx.Done()
		b.removeSubscriber(organizationID, subscriber)
	}()

	return subscriber, nil
}

func (b *memoryBroadcaster) removeSubscriber(organizationID uuid.UUID, subscriber chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	room, ok := b.subscribers[organizationID]
	if !ok {
		return
	}
	if _, subscribed := room[subscriber]; !subscribed {
		return
	}

	delete(room, subscriber)
	if len(room) == 0 {
		delete(b.subscribers, organizationID)
	}
	close(subscriber)
}

func (b *memoryBroadcaster) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}
	b.closed = true

	for organizationID, room := range b.subscribers {
		for subscriber := range room {
			close(subscriber)
		}
		delete(b.subscribers, organizationID)
	}

	return nil
}
