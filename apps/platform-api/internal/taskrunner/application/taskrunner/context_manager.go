package taskrunner

import (
	"context"
	"sync"

	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
	"github.com/google/uuid"
)

type ContextManager struct {
	contexts map[uuid.UUID]context.CancelFunc
	mu       sync.RWMutex
}

func NewContextManager() taskRunnerDomainTaskRunner.ContextManager {
	return &ContextManager{
		contexts: make(map[uuid.UUID]context.CancelFunc),
	}
}

func (cm *ContextManager) CreateContext(taskID uuid.UUID, parentCtx context.Context) context.Context {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cancel, exists := cm.contexts[taskID]; exists {
		cancel()
	}

	ctx, cancel := context.WithCancel(parentCtx)
	cm.contexts[taskID] = cancel
	return ctx
}

func (cm *ContextManager) CancelContext(taskID uuid.UUID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cancel, exists := cm.contexts[taskID]; exists {
		cancel()
		delete(cm.contexts, taskID)
	}
}

func (cm *ContextManager) RemoveContext(taskID uuid.UUID) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.contexts, taskID)
}

func (cm *ContextManager) CancelAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for taskID, cancel := range cm.contexts {
		cancel()
		delete(cm.contexts, taskID)
	}
}
