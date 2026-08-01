package deletetrigger

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
)

type Handler struct {
	triggerRepository  triggersDomainTriggers.TriggerRepository
	transactionManager database.TransactionManager
}

func New(
	triggerRepository triggersDomainTriggers.TriggerRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		triggerRepository:  triggerRepository,
		transactionManager: transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *DeleteTriggerCommand) (*DeleteTriggerResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		trigger, err := h.triggerRepository.GetByIDAndOrganizationID(txCtx, command.TriggerID, command.OrganizationID)
		if err != nil {
			return triggersDomainTriggers.ErrTriggerNotFound
		}

		deletedTrigger, err := trigger.Delete()
		if err != nil {
			return err
		}

		return h.triggerRepository.Delete(txCtx, deletedTrigger)
	})

	if err != nil {
		return nil, err
	}

	return &DeleteTriggerResponse{}, nil
}
