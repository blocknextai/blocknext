package rerunfailed

import (
	"context"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
)

type RerunFailedHandler struct {
	taskService taskRunnerDomainTaskRunner.TaskService
}

func NewRerunFailedHandler(
	taskService taskRunnerDomainTaskRunner.TaskService,
) *RerunFailedHandler {
	return &RerunFailedHandler{
		taskService: taskService,
	}
}

func (h *RerunFailedHandler) Handle(ctx context.Context, command *RerunFailedCommand) (*RerunFailedResponse, error) {
	newTask, err := h.taskService.RerunFailed(
		ctx,
		new(command.TriggeredByUserID),
		command.OrganizationID,
		command.ID,
		taskRunnerDomainTask.TaskTriggerTypeRerunFailed,
	)

	if err != nil {
		return nil, err
	}

	return &RerunFailedResponse{
		ID: newTask.ID,
	}, nil
}
