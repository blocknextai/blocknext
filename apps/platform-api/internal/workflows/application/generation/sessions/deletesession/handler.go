package deletesession

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	generationDomainMessages "github.com/blocknextai/platform-api/internal/workflows/domain/generation/messages"
	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
)

type Handler struct {
	sessionRepository  generationDomainSessions.SessionRepository
	messageRepository  generationDomainMessages.MessageRepository
	transactionManager database.TransactionManager
}

func New(
	sessionRepository generationDomainSessions.SessionRepository,
	messageRepository generationDomainMessages.MessageRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		sessionRepository:  sessionRepository,
		messageRepository:  messageRepository,
		transactionManager: transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, request *DeleteSessionCommand) (*DeleteSessionResponse, error) {
	session, err := h.sessionRepository.GetByIDAndOrganizationID(ctx, request.SessionID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	deletedSession, err := session.Delete()
	if err != nil {
		return nil, err
	}

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		if err := h.messageRepository.DeleteBySessionID(txCtx, deletedSession.ID); err != nil {
			return err
		}

		if err := h.sessionRepository.Delete(txCtx, deletedSession); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &DeleteSessionResponse{
		Success: true,
	}, nil
}
