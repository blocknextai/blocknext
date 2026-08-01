package logout

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

func (h *Handler) Handle(ctx context.Context, command *LogoutCommand) (*LogoutResponse, error) {
	if err := h.sessionService.RevokeSession(ctx, command.SessionID); err != nil {
		return nil, err
	}

	return &LogoutResponse{}, nil
}
