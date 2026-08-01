package revokesession

import (
	"context"

	accountApplicationSessions "github.com/blocknextai/platform-api/internal/account/application/sessions"
	accountDomainSessions "github.com/blocknextai/platform-api/internal/account/domain/sessions"
)

type Handler struct {
	sessionService    accountApplicationSessions.SessionService
	sessionRepository accountDomainSessions.SessionRepository
}

func New(
	sessionService accountApplicationSessions.SessionService,
	sessionRepository accountDomainSessions.SessionRepository,
) *Handler {
	return &Handler{
		sessionService:    sessionService,
		sessionRepository: sessionRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, command *RevokeSessionCommand) (*RevokeSessionResponse, error) {
	if _, err := h.sessionRepository.GetByIDAndUserID(ctx, command.SessionID, command.UserID); err != nil {
		return nil, err
	}

	if err := h.sessionService.RevokeSession(ctx, command.SessionID); err != nil {
		return nil, err
	}

	return &RevokeSessionResponse{}, nil
}
