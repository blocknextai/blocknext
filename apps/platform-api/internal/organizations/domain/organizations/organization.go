package organizations

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/uuid"
)

type Organization struct {
	database.BaseEntity

	Title       string
	Description *string
	IsVerified  bool
}

func New(
	title string,
	description *string,
	isVerified bool,
) (*Organization, error) {
	utcNow := time.Now().UTC()

	organization := &Organization{
		BaseEntity: database.BaseEntity{
			ID:        uuid.NewV7(),
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		Title:       title,
		Description: description,
		IsVerified:  isVerified,
	}

	return organization.validateThenReturn()
}

func (o *Organization) Update(
	title string,
	description *string,
) (*Organization, error) {
	o.UpdatedAt = time.Now().UTC()

	o.Title = title
	o.Description = description

	return o.validateThenReturn()
}

func (o *Organization) Delete() (*Organization, error) {
	utcNow := time.Now().UTC()

	o.UpdatedAt = utcNow
	o.DeletedAt = new(utcNow)
	return o.validateThenReturn()
}

func (o *Organization) validateThenReturn() (*Organization, error) {
	if strings.TrimSpace(o.Title) == "" {
		return nil, ErrTitleIsRequired
	}
	return o, nil
}
