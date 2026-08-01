package executors

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
)

type ExecutorService interface {
	GetByID(id string) (executors.ExecutorManager, bool)
}

type executorService struct{}

func NewExecutorService() ExecutorService {
	return &executorService{}
}

func (s *executorService) GetByID(id string) (executors.ExecutorManager, bool) {
	return executors.GetExecutor(id)
}
