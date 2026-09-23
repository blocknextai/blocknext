package tasks

import (
	"context"

	"github.com/google/uuid"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/task"
)

type RerunFailedCommand struct {
	TriggeredByUserID uuid.UUID
	OrganizationID    uuid.UUID
	ID                uuid.UUID
}

type RerunFailedResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Service) RerunFailed(ctx context.Context, command *RerunFailedCommand) (*RerunFailedResponse, error) {
	newTask, err := s.taskService.RerunFailed(
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
