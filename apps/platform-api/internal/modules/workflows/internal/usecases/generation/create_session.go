package generation

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	generationDomainSessions "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/generation/sessions"
)

type CreateSessionCommand struct {
	OrganizationID *uuid.UUID
	UserID         *uuid.UUID
	Title          string
}

func (c *CreateSessionCommand) Validate() error {
	if strings.TrimSpace(c.Title) == "" {
		return generationDomainSessions.ErrTitleIsRequired
	}

	return nil
}

type CreateSessionResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *Service) CreateSession(ctx context.Context, request *CreateSessionCommand) (*CreateSessionResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	session, err := generationDomainSessions.New(
		request.OrganizationID,
		request.UserID,
		request.Title,
	)
	if err != nil {
		return nil, err
	}

	if err := s.sessionRepository.Create(ctx, session); err != nil {
		return nil, err
	}

	return &CreateSessionResponse{
		ID:        session.ID,
		Title:     session.Title,
		CreatedAt: session.CreatedAt,
	}, nil
}
