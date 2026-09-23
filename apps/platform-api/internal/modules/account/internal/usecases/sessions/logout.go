package sessions

import (
	"context"

	"github.com/google/uuid"
)

type LogoutCommand struct {
	SessionID uuid.UUID
}

type LogoutResponse struct{}

func (s *Service) Logout(ctx context.Context, command *LogoutCommand) (*LogoutResponse, error) {
	if err := s.sessionService.RevokeSession(ctx, command.SessionID); err != nil {
		return nil, err
	}

	return &LogoutResponse{}, nil
}
