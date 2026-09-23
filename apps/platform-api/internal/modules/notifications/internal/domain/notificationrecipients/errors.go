package notificationrecipients

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrNotificationRecipientNotFound = apperror.NotFound("notification recipient not found")
	ErrInvalidNotificationID         = apperror.Validation("invalid notification id")
	ErrInvalidUserID                 = apperror.Validation("invalid user id")
)
