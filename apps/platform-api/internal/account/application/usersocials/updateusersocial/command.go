package updateusersocial

import (
	"github.com/google/uuid"
)

type UpdateUserSocialItem struct {
	URL string
}

type UpdateUserSocialCommand struct {
	UserID uuid.UUID
	Items  []UpdateUserSocialItem
}
