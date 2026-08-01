package createsession

import (
	"context"

	generationDomainSessions "github.com/blocknextai/platform-api/internal/workflows/domain/generation/sessions"
)

type Handler struct {
	sessionRepository generationDomainSessions.SessionRepository
}

func New(
	sessionRepository generationDomainSessions.SessionRepository,
) *Handler {
	return &Handler{
		sessionRepository: sessionRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *CreateSessionCommand) (*CreateSessionResponse, error) {
	session, err := generationDomainSessions.New(
		request.OrganizationID,
		request.UserID,
		request.Title,
	)
	if err != nil {
		return nil, err
	}

	if err := h.sessionRepository.Create(ctx, session); err != nil {
		return nil, err
	}

	return &CreateSessionResponse{
		ID:        session.ID,
		Title:     session.Title,
		CreatedAt: session.CreatedAt,
	}, nil
}
