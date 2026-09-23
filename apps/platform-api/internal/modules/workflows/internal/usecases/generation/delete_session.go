package generation

import (
	"context"

	"github.com/google/uuid"
)

type DeleteSessionCommand struct {
	OrganizationID uuid.UUID
	SessionID      uuid.UUID
}

type DeleteSessionResponse struct {
	Success bool `json:"success"`
}

func (s *Service) DeleteSession(ctx context.Context, request *DeleteSessionCommand) (*DeleteSessionResponse, error) {
	session, err := s.sessionRepository.GetByIDAndOrganizationID(ctx, request.SessionID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	deletedSession, err := session.Delete()
	if err != nil {
		return nil, err
	}

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		if err := s.messageRepository.DeleteBySessionID(txCtx, deletedSession.ID); err != nil {
			return err
		}

		if err := s.sessionRepository.Delete(txCtx, deletedSession); err != nil {
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
