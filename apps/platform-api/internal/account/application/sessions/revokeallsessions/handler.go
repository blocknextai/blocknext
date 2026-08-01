package revokeallsessions

import (
	"context"

	accountApplicationSessions "github.com/blocknextai/platform-api/internal/account/application/sessions"
)

type Handler struct {
	sessionService accountApplicationSessions.SessionService
}

func New(
	sessionService accountApplicationSessions.SessionService,
) *Handler {
	return &Handler{
		sessionService: sessionService,
	}
}

func (h *Handler) Handle(ctx context.Context, command *RevokeAllSessionsCommand) (*RevokeAllSessionsResponse, error) {
	if _, err := h.sessionService.RevokeAllUserSessions(ctx, command.UserID); err != nil {
		return nil, err
	}

	return &RevokeAllSessionsResponse{}, nil
}
