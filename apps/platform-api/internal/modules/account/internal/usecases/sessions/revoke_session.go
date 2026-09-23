package sessions

import (
	"context"

	"github.com/google/uuid"
)

type RevokeSessionCommand struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
}

type RevokeSessionResponse struct{}

func (s *Service) RevokeSession(ctx context.Context, command *RevokeSessionCommand) (*RevokeSessionResponse, error) {
	if _, err := s.sessionRepository.GetByIDAndUserID(ctx, command.SessionID, command.UserID); err != nil {
		return nil, err
	}

	if err := s.sessionService.RevokeSession(ctx, command.SessionID); err != nil {
		return nil, err
	}

	return &RevokeSessionResponse{}, nil
}
