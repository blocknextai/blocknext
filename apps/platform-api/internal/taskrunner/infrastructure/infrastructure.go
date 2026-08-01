package infrastructure

import (
	"github.com/blocknextai/platform-api/internal/taskrunner/application/canceltask"
	"github.com/blocknextai/platform-api/internal/taskrunner/application/contextresolver"
	"github.com/blocknextai/platform-api/internal/taskrunner/application/rerunall"
	"github.com/blocknextai/platform-api/internal/taskrunner/application/rerunfailed"
	"github.com/blocknextai/platform-api/internal/taskrunner/application/triggertask"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
)

type Handlers struct {
	TriggerTask *triggertask.TriggerTaskHandler
	CancelTask  *canceltask.CancelTaskHandler
	RerunAll    *rerunall.RerunAllHandler
	RerunFailed *rerunfailed.RerunFailedHandler
}

type RegisterInfrastructureDeps struct {
	TaskService     taskRunnerDomainTaskRunner.TaskService
	ContextResolver contextresolver.ContextResolver
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		TriggerTask: triggertask.NewTriggerTaskHandler(deps.TaskService, deps.ContextResolver),
		CancelTask:  canceltask.NewCancelTaskHandler(deps.TaskService),
		RerunAll:    rerunall.NewRerunAllHandler(deps.TaskService),
		RerunFailed: rerunfailed.NewRerunFailedHandler(deps.TaskService),
	}
}
