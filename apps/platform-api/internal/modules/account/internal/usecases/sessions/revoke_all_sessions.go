package sessions

import (
	"context"

	"github.com/google/uuid"
)

type RevokeAllSessionsCommand struct {
	UserID uuid.UUID
}

type RevokeAllSessionsResponse struct{}

func (s *Service) RevokeAllSessions(ctx context.Context, command *RevokeAllSessionsCommand) (*RevokeAllSessionsResponse, error) {
	if _, err := s.sessionService.RevokeAllUserSessions(ctx, command.UserID); err != nil {
		return nil, err
	}

	return &RevokeAllSessionsResponse{}, nil
}
