package generation

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	generationDomainSessions "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/generation/sessions"
)

type UpdateSessionCommand struct {
	OrganizationID uuid.UUID
	SessionID      uuid.UUID
	Title          string
}

func (c *UpdateSessionCommand) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return generationDomainSessions.ErrTitleIsRequired
	}

	return nil
}

type UpdateSessionResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Service) UpdateSession(ctx context.Context, request *UpdateSessionCommand) (*UpdateSessionResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	session, err := s.sessionRepository.GetByIDAndOrganizationID(ctx, request.SessionID, request.OrganizationID)
	if err != nil {
		return nil, err
	}

	session, err = session.Update(request.Title)
	if err != nil {
		return nil, err
	}

	if err := s.sessionRepository.Update(ctx, session); err != nil {
		return nil, err
	}

	return &UpdateSessionResponse{
		ID:        session.ID,
		Title:     session.Title,
		UpdatedAt: session.UpdatedAt,
	}, nil
}
