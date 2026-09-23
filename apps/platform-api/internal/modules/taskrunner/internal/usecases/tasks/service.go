package tasks

import (
	taskRunnerApplicationContextresolver "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/contextresolver"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/taskrunner"
)

type Service struct {
	taskService     taskRunnerDomainTaskRunner.TaskService
	contextResolver taskRunnerApplicationContextresolver.ContextResolver
}

func NewService(
	taskService taskRunnerDomainTaskRunner.TaskService,
	contextResolver taskRunnerApplicationContextresolver.ContextResolver,
) *Service {
	return &Service{
		taskService:     taskService,
		contextResolver: contextResolver,
	}
}
