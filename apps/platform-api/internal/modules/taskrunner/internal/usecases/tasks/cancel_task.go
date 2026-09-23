package tasks

import (
	"context"

	"github.com/google/uuid"
)

type CancelTaskCommand struct {
	TriggeredByUserID uuid.UUID
	OrganizationID    uuid.UUID
	ID                uuid.UUID
}

type CancelTaskResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Service) CancelTask(ctx context.Context, command *CancelTaskCommand) (*CancelTaskResponse, error) {
	err := s.taskService.CancelTask(
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
