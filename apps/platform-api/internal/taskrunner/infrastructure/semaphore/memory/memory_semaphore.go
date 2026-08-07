package memory

import (
	"context"
	"sync"
	"time"

	"github.com/blocknextai/platform-api/internal/taskrunner/infrastructure/semaphore"
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

func (sm *MemorySemaphore) AcquireSemaphore(ctx context.Context, organizationID uuid.UUID, taskID uuid.UUID, maxConcurrentExecutions int64) (chan struct{}, error) {
	return semaphore.AcquireWithBackoff(ctx, func(_ context.Context) (bool, error) {
		return sm.tryAcquire(organizationID, taskID, maxConcurrentExecutions), nil
	})
}

func (sm *MemorySemaphore) tryAcquire(organizationID uuid.UUID, taskID uuid.UUID, maxConcurrentExecutions int64) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	room := sm.pruneExpired(organizationID)

	if _, held := room[taskID]; held {
		room[taskID] = time.Now().Add(sm.ttl)
		return true
	}

	if int64(len(room)) >= maxConcurrentExecutions {
		return false
	}

	room[taskID] = time.Now().Add(sm.ttl)
	sm.holders[organizationID] = room

	return true
}

func (sm *MemorySemaphore) ReleaseSemaphore(_ context.Context, organizationID uuid.UUID, taskID uuid.UUID, token chan struct{}) error {
	semaphore.Drain(token)

	sm.mu.Lock()
	defer sm.mu.Unlock()

	room, ok := sm.holders[organizationID]
	if !ok {
		return nil
	}

	delete(room, taskID)
	if len(room) == 0 {
		delete(sm.holders, organizationID)
	}

	return nil
}

func (sm *MemorySemaphore) HeartbeatSemaphore(_ context.Context, organizationID uuid.UUID, taskID uuid.UUID) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	room := sm.pruneExpired(organizationID)
	if _, held := room[taskID]; !held {
		return nil
	}

	room[taskID] = time.Now().Add(sm.ttl)
	sm.holders[organizationID] = room

	return nil
}

func (sm *MemorySemaphore) pruneExpired(organizationID uuid.UUID) map[uuid.UUID]time.Time {
	room, ok := sm.holders[organizationID]
	if !ok {
		return make(map[uuid.UUID]time.Time)
	}

	now := time.Now()
	for taskID, expiresAt := range room {
		if now.After(expiresAt) {
			delete(room, taskID)
		}
	}

	return room
}
