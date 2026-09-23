package usecases

import (
	"github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/application/contextresolver"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/taskrunner"
	tasksUseCases "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/usecases/tasks"
)

type Services struct {
	Tasks *tasksUseCases.Service
}

type ServiceDependencies struct {
	TaskService     taskRunnerDomainTaskRunner.TaskService
	ContextResolver contextresolver.ContextResolver
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		Tasks: tasksUseCases.NewService(deps.TaskService, deps.ContextResolver),
	}
}
