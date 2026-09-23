package triggers

import (
	"context"

	"github.com/google/uuid"

	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
)

type DeleteTriggerCommand struct {
	OrganizationID uuid.UUID
	TriggerID      uuid.UUID
}

type DeleteTriggerResponse struct{}

func (s *Service) DeleteTrigger(ctx context.Context, command *DeleteTriggerCommand) (*DeleteTriggerResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		trigger, err := s.triggerRepository.GetByIDAndOrganizationID(txCtx, command.TriggerID, command.OrganizationID)
		if err != nil {
			return triggersDomainTriggers.ErrTriggerNotFound
		}

		deletedTrigger, err := trigger.Delete()
		if err != nil {
			return err
		}

		return s.triggerRepository.Delete(txCtx, deletedTrigger)
	})

	if err != nil {
		return nil, err
	}

	return &DeleteTriggerResponse{}, nil
}
