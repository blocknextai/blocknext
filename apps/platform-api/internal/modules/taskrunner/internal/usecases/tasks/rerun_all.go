package tasks

import (
	"context"

	"github.com/google/uuid"

	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/modules/taskrunner/internal/domain/task"
)

type RerunAllCommand struct {
	TriggeredByUserID uuid.UUID
	OrganizationID    uuid.UUID
	ID                uuid.UUID
}

type RerunAllResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Service) RerunAll(ctx context.Context, command *RerunAllCommand) (*RerunAllResponse, error) {
	newTask, err := s.taskService.RerunAll(
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
