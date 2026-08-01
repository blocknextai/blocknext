package canceltask

import (
	"context"

	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
)

type CancelTaskHandler struct {
	taskService taskRunnerDomainTaskRunner.TaskService
}

func NewCancelTaskHandler(
	taskService taskRunnerDomainTaskRunner.TaskService,
) *CancelTaskHandler {
	return &CancelTaskHandler{
		taskService: taskService,
	}
}

func (h *CancelTaskHandler) Handle(ctx context.Context, command *CancelTaskCommand) (*CancelTaskResponse, error) {
	err := h.taskService.CancelTask(
		ctx,
		new(command.TriggeredByUserID),
		command.OrganizationID,
		command.ID,
	)

	if err != nil {
		return nil, err
	}

	return &CancelTaskResponse{
		ID: command.ID,
	}, nil
}
