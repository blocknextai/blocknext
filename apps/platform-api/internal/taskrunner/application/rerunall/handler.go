package rerunall

import (
	"context"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
)

type RerunAllHandler struct {
	taskService taskRunnerDomainTaskRunner.TaskService
}

func NewRerunAllHandler(
	taskService taskRunnerDomainTaskRunner.TaskService,
) *RerunAllHandler {
	return &RerunAllHandler{
		taskService: taskService,
	}
}

func (h *RerunAllHandler) Handle(ctx context.Context, command *RerunAllCommand) (*RerunAllResponse, error) {
	newTask, err := h.taskService.RerunAll(
		ctx,
		new(command.TriggeredByUserID),
		command.OrganizationID,
		command.ID,
		taskRunnerDomainTask.TaskTriggerTypeRerunAll,
	)

	if err != nil {
		return nil, err
	}

	return &RerunAllResponse{
		ID: newTask.ID,
	}, nil
}
