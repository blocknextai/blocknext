package users

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/rbac"
	"github.com/blocknextai/go-packages/uuid"
)

type User struct {
	database.BaseEntity

	Role       string
	IsVerified bool
	IsBanned   bool
}

func NewUser() (*User, error) {
	utcNow := time.Now().UTC()

	user := &User{
		ID:        uuid.NewV7(),
		CreatedAt: utcNow,
		UpdatedAt: utcNow,
		DeletedAt: nil,
		Role:      rbac.GlobalAdminRole.Name,
	}

	return user.validateThenReturn()
}

func (u *User) validateThenReturn() (*User, error) {
	if strings.TrimSpace(u.Role) == "" {
		return nil, ErrRoleIsRequired
	}

	if !rbac.IsValidUserRole(u.Role) {
		return nil, ErrInvalidRole
	}

	return u, nil
}
