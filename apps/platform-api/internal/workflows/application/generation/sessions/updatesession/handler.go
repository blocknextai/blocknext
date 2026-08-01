package updatesession

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

func (h *Handler) Handle(ctx context.Context, request *UpdateSessionCommand) (*UpdateSessionResponse, error) {
	session, err := h.sessionRepository.GetByIDAndOrganizationID(ctx, request.SessionID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	session, err = session.Update(request.Title)
	if err != nil {
		return nil, err
	}

	if err := h.sessionRepository.Update(ctx, session); err != nil {
		return nil, err
	}

	return &UpdateSessionResponse{
		ID:        session.ID,
		Title:     session.Title,
		UpdatedAt: session.UpdatedAt,
	}, nil
}
