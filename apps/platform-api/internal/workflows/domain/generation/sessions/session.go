package sessions

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type GenerationSession struct {
	database.BaseEntity

	OrganizationID *uuid.UUID
	UserID         *uuid.UUID
	Title          string
}

func New(
	organizationID *uuid.UUID,
	userID *uuid.UUID,
	title string,
) (*GenerationSession, error) {
	utcNow := time.Now().UTC()

	session := &GenerationSession{
		ID:             bnuuid.NewV7(),
		CreatedAt:      utcNow,
		UpdatedAt:      utcNow,
		DeletedAt:      nil,
		OrganizationID: organizationID,
		UserID:         userID,
		Title:          title,
	}

	return session.validateThenReturn()
}

func (s *GenerationSession) Update(title string) (*GenerationSession, error) {
	s.UpdatedAt = time.Now().UTC()

	s.Title = title

	return s.validateThenReturn()
}

func (s *GenerationSession) Delete() (*GenerationSession, error) {
	utcNow := time.Now().UTC()

	s.UpdatedAt = utcNow
	s.DeletedAt = new(utcNow)
	return s.validateThenReturn()
}

func (s *GenerationSession) validateThenReturn() (*GenerationSession, error) {
	if strings.TrimSpace(s.Title) == "" {
		return nil, ErrTitleIsRequired
	}
	return s, nil
}
