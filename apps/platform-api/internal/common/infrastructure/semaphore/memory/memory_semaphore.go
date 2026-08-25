package memory

import (
	"context"
	"sync"
	"time"

	"github.com/blocknextai/platform-api/internal/common/infrastructure/semaphore/token"
	"github.com/google/uuid"
)

type MemorySemaphore struct {
	mu      sync.Mutex
	holders map[uuid.UUID]map[uuid.UUID]time.Time
	ttl     time.Duration
}

func New(ttl time.Duration) *MemorySemaphore {
	return &MemorySemaphore{
		holders: make(map[uuid.UUID]map[uuid.UUID]time.Time),
		ttl:     ttl,
	}
}

func (sm *MemorySemaphore) Ping(_ context.Context) error {
	return nil
}

func (sm *MemorySemaphore) AcquireSemaphore(ctx context.Context, organizationID uuid.UUID, holderID uuid.UUID, maxConcurrentExecutions int64) (chan struct{}, error) {
	return token.AcquireWithBackoff(ctx, func(_ context.Context) (bool, error) {
		return sm.tryAcquire(organizationID, holderID, maxConcurrentExecutions), nil
	})
}

func (sm *MemorySemaphore) tryAcquire(organizationID uuid.UUID, holderID uuid.UUID, maxConcurrentExecutions int64) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	room := sm.pruneExpired(organizationID)

	if _, held := room[holderID]; held {
		room[holderID] = time.Now().Add(sm.ttl)
		return true
	}

	if int64(len(room)) >= maxConcurrentExecutions {
		return false
	}

	room[holderID] = time.Now().Add(sm.ttl)
	sm.holders[organizationID] = room

	return true
}

func (sm *MemorySemaphore) ReleaseSemaphore(_ context.Context, organizationID uuid.UUID, holderID uuid.UUID, semaphore chan struct{}) error {
	token.Drain(semaphore)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	room, ok := sm.holders[organizationID]
	if !ok {
		return nil
	}

	delete(room, holderID)
	if len(room) == 0 {
		delete(sm.holders, organizationID)
	}

	return nil
}

func (sm *MemorySemaphore) HeartbeatSemaphore(_ context.Context, organizationID uuid.UUID, holderID uuid.UUID) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	room := sm.pruneExpired(organizationID)
	if _, held := room[holderID]; !held {
		return nil
	}

	room[holderID] = time.Now().Add(sm.ttl)
	sm.holders[organizationID] = room

	return nil
}

func (sm *MemorySemaphore) pruneExpired(organizationID uuid.UUID) map[uuid.UUID]time.Time {
	room, ok := sm.holders[organizationID]
	if !ok {
		return make(map[uuid.UUID]time.Time)
	}

	now := time.Now()
	for holderID, expiresAt := range room {
		if now.After(expiresAt) {
			delete(room, holderID)
		}
	}

	return room
}
