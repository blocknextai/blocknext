package usersocials

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type UserSocial struct {
	database.BaseEntity

	UserID    uuid.UUID
	Platform  string
	URL       string
	SortOrder int
}

func NewUserSocial(
	userID uuid.UUID,
	platform string,
	url string,
	sortOrder int,
) (*UserSocial, error) {
	utcNow := time.Now().UTC()

	social := &UserSocial{
		ID:        bnuuid.NewV7(),
		CreatedAt: utcNow,
		UpdatedAt: utcNow,
		DeletedAt: nil,
		UserID:    userID,
		Platform:  platform,
		URL:       url,
		SortOrder: sortOrder,
	}

	return social.validateThenReturn()
}

func (us *UserSocial) Delete() (*UserSocial, error) {
	utcNow := time.Now().UTC()

	us.UpdatedAt = utcNow
	us.DeletedAt = new(utcNow)
	return us.validateThenReturn()
}

func (us *UserSocial) validateThenReturn() (*UserSocial, error) {
	if us.UserID == uuid.Nil {
		return nil, ErrUserIDIsRequired
	}
	if strings.TrimSpace(us.Platform) == "" {
		return nil, ErrPlatformIsRequired
	}
	if strings.TrimSpace(us.URL) == "" {
		return nil, ErrURLIsRequired
	}

	return us, nil
}
