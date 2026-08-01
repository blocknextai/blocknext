package executors

import (
	"sync"
)

var (
	once     sync.Once
	instance *ExecutorRegistry
)

type ExecutorRegistry struct {
	executorsMap map[string]ExecutorManager
}

func GetRegistry() *ExecutorRegistry {
	once.Do(func() {
		instance = &ExecutorRegistry{
			executorsMap: make(map[string]ExecutorManager),
		}
	})
	return instance
}

func (r *ExecutorRegistry) RegisterExecutor(executor ExecutorManager) {
	r.executorsMap[executor.GetID()] = executor
}

func (r *ExecutorRegistry) GetExecutor(id string) (ExecutorManager, bool) {
	executor, ok := r.executorsMap[id]
	return executor, ok
}

func RegisterExecutor(executor ExecutorManager) {
	if executor.GetDisabled() {
		return
	}

	GetRegistry().RegisterExecutor(executor)
}

func GetExecutor(id string) (ExecutorManager, bool) {
	return GetRegistry().GetExecutor(id)
}
